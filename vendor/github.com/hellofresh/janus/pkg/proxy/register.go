package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/hellofresh/janus/pkg/router"
	log "github.com/sirupsen/logrus"
)

const (
	methodAll = "ALL"
)

// Register handles the register of proxies into the chosen router.
// It also handles the conversion from a proxy to an http.HandlerFunc
type Register struct {
	Router router.Router
	params Params
}

// NewRegister creates a new instance of Register
func NewRegister(router router.Router, params Params) *Register {
	return &Register{router, params}
}

// AddMany registers many proxies at once
func (p *Register) AddMany(routes []*Route) error {
	for _, r := range routes {
		err := p.Add(r)
		if nil != err {
			return err
		}
	}

	return nil
}

// Add register a new route
func (p *Register) Add(route *Route) error {
	definition := route.Proxy

	p.params.Outbound = route.Outbound
	p.params.InsecureSkipVerify = definition.InsecureSkipVerify
	handler := &httputil.ReverseProxy{
		Director:  p.createDirector(definition),
		Transport: NewTransportWithParams(p.params),
	}

	matcher := router.NewListenPathMatcher()
	if matcher.Match(definition.ListenPath) {
		p.doRegister(matcher.Extract(definition.ListenPath), handler.ServeHTTP, definition.Methods, route.Inbound)
	}

	p.doRegister(definition.ListenPath, handler.ServeHTTP, definition.Methods, route.Inbound)
	return nil
}

func (p *Register) createDirector(proxyDefinition *Definition) func(req *http.Request) {
	return func(req *http.Request) {
		var target *url.URL
		var err error

		// TODO: find better solution
		// maybe create "proxyDefinition.Upstreams.Targets every time",
		// but currently we have several points of definition creation
		if proxyDefinition.Upstreams != nil && proxyDefinition.Upstreams.Targets != nil && len(proxyDefinition.Upstreams.Targets) > 0 {
			log.WithField("balancing_alg", proxyDefinition.Upstreams.Balancing).Debug("Using a load balancing algorithm")
			balancer, err := NewBalancer(proxyDefinition.Upstreams.Balancing)
			if err != nil {
				log.WithError(err).Error("Could not create a balancer")
				return
			}

			upstream, err := balancer.Elect(proxyDefinition.Upstreams.Targets)
			if err != nil {
				log.WithError(err).Error("Could not elect one upstream")
				return
			}

			log.WithField("target", upstream.Target).Debug("Elected Target")
			target, err = url.Parse(upstream.Target)
			if err != nil {
				log.WithError(err).Error("Could not parse the target URL")
				return
			}
		} else {
			log.Warn("The upstream URL is deprecated. Use Upstreams instead")
			target, err = url.Parse(proxyDefinition.UpstreamURL)
			if err != nil {
				log.WithError(err).Error("Could not parse the target URL")
				return
			}
		}

		targetQuery := target.RawQuery

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		path := target.Path

		if proxyDefinition.AppendPath {
			log.Debug("Appending listen path to the target url")
			path = singleJoiningSlash(target.Path, req.URL.Path)
		}

		if proxyDefinition.StripPath {
			path = singleJoiningSlash(target.Path, req.URL.Path)
			matcher := router.NewListenPathMatcher()
			listenPath := matcher.Extract(proxyDefinition.ListenPath)

			log.WithField("listen_path", listenPath).Debug("Stripping listen path")
			path = strings.Replace(path, listenPath, "", 1)
			if !strings.HasSuffix(target.Path, "/") && strings.HasSuffix(path, "/") {
				path = path[:len(path)-1]
			}
		}

		log.WithField("path", path).Debug("Upstream Path")
		req.URL.Path = path

		// This is very important to avoid problems with ssl verification for the HOST header
		if !proxyDefinition.PreserveHost {
			log.Debug("Preserving the host header")
			req.Host = target.Host
		}

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}
}

func (p *Register) doRegister(listenPath string, handler http.HandlerFunc, methods []string, handlers InChain) {
	log.WithFields(log.Fields{
		"listen_path": listenPath,
	}).Debug("Registering a route")

	if strings.Index(listenPath, "/") != 0 {
		log.WithField("listen_path", listenPath).
			Error("Route listen path must begin with '/'.Skipping invalid route.")
	} else {
		for _, method := range methods {
			if strings.ToUpper(method) == methodAll {
				p.Router.Any(listenPath, handler, handlers...)
			} else {
				p.Router.Handle(strings.ToUpper(method), listenPath, handler, handlers...)
			}
		}
	}
}

func cleanSlashes(a string) string {
	endSlash := strings.HasSuffix(a, "//")
	startSlash := strings.HasPrefix(a, "//")

	if startSlash {
		a = "/" + strings.TrimPrefix(a, "//")
	}

	if endSlash {
		a = strings.TrimSuffix(a, "//") + "/"
	}

	return a
}

func singleJoiningSlash(a, b string) string {
	a = cleanSlashes(a)
	b = cleanSlashes(b)

	aSlash := strings.HasSuffix(a, "/")
	bSlash := strings.HasPrefix(b, "/")

	switch {
	case aSlash && bSlash:
		return a + b[1:]
	case !aSlash && !bSlash:
		if len(b) > 0 {
			return a + "/" + b
		}
		return a
	}
	return a + b
}

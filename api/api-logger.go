package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/hellofresh/janus/pkg/response"

	"github.com/sirupsen/logrus"
)

var api API

//LogHandler log any resquest incoming to API
func LogHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		api.Log = logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		})

		api.Log.Debug("Request's starting")
		// originalURL := &url.URL{}
		// *originalURL := *r.URL

		var (
			lock         sync.Mutex
			responseCode int
		)
		hooks := response.Hooks{
			WriteHeader: func(next response.WriteHeaderFunc) response.WriteHeaderFunc {
				return func(code int) {
					next(code)
					lock.Lock()
					defer lock.Unlock()

					responseCode = code
				}
			},
		}

		fields := logrus.Fields{
			"methode":     r.Method,
			"host":        r.Host,
			"request":     r.RequestURI,
			"remote-addr": r.RemoteAddr,
			"user-agent":  r.UserAgent(),
			"referer":     r.Referer(),
		}

		h.ServeHTTP(response.Wrap(w, hooks), r)

		fields["code"] = responseCode
		fields["duration"] = int(time.Now().Sub(startTime) / time.Millisecond)

		api.Log.WithFields(fields).Info("completed handling request")
	})
}

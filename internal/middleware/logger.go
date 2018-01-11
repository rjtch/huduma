package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/huduma/internal/router"
)

//Logger logs requests by writing some infos in the format
//rqID: (200) GET /foo -> ADDR (latency)
func Logger(h router.Handler) router.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

		//we use to log all request by insering some informations about them.
		v := ctx.Value(router.KeyV).(*router.Value)
		start := v.Now

		h(ctx, w, r, params)
		log.Printf("%s : (%d) : (%s) %s -> %s (%s)",
			v.RqID,
			v.State,
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			time.Since(start),
		)

		return nil
	}
}

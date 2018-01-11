package router

import (
	"context"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"

	"github.com/pborman/uuid"
)

//Handler is a type http.handler that handles an http request into our logic.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error

//ReqIDHeader is the ID to be added to the header of outgoing requests.
const ReqIDHeader = "X-Trace-ID"

//ctxtKey represents the value of the context key.
type ctxtKey int

//KeyV is how the value is stored/retreived.
const KeyV ctxtKey = 1

//Value represents the state for each request
type Value struct {
	RqID  string
	Now   time.Time
	State int
}

//Huduma is the struture of our router and entrypoint of any request that will log and add context value to
//any incoming request
type Huduma struct {
	*httptreemux.TreeMux
	mw []Middleware
}

//NewHuduma creates a new Huduma's instance
func NewHuduma(mw ...Middleware) *Huduma {

	router := httptreemux.New()
	return &Huduma{
		TreeMux: router,
		mw:      mw,
	}
}

//HandleWrapper is used to mount http verbs and pair together for convenience routing
//verb can be any http methods
func (h *Huduma) HandleWrapper(verb, path string, handler Handler, m ...Middleware) {

	//handler wraps up the whole applicaltion and call the first function of each middleware
	//which will return a function of type handler.
	handler = wrapMiddleware(wrapMiddleware(handler, m), h.mw)

	a := func(w http.ResponseWriter, r *http.Request, params map[string]string) {

		//check the request for an existing traceID. If it doesn't exist, generate a new one.
		rqID := r.Header.Get(ReqIDHeader)
		if rqID == "" {
			rqID = uuid.New()
		}

		//set the context with the required value to process the request
		v := Value{
			RqID: rqID,
			Now:  time.Now(),
		}
		ctx := context.WithValue(r.Context(), KeyV, &v)

		//Set the request id to outgoing requests before any other header to ensure that the request id
		//is always added to the request regardless of any error occuring or not.
		w.Header().Set(ReqIDHeader, v.RqID)

		handler(ctx, w, r, params)
	}

	//Add this handler to the specified verb and route
	h.TreeMux.Handle(verb, path, a)
}

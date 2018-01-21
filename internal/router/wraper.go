package router

// Middleware is an http handler that has one input argument and one output argument of type http.handler
//It is used to remove boilerplate or other concerns not direct to any given Handler.
type Middleware func(Handler) Handler

//wrapMiddleware is the middleware wrapper
//It wraps many handler to be log, trace, validate or to which response hearders will be add.
func wrapMiddleware(handler Handler, m []Middleware) Handler {

	for i := len(m) - 1; i >= 0; i-- {
		if m[i] != nil {
			handler = m[i](handler)
		}
	}
	return handler
}

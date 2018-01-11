package middleware

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/pkg/errors"

	"github.com/huduma/internal/router"
)

/*
DON'T JUST CHECK ERRORS HANDLE THEM GRACEFULLY. GO PROVERB.
*/

//ErrorHandler is used to catch and to respond errors in middleware.
func ErrorHandler(h router.Handler) router.Handler {

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

		v := ctx.Value(router.KeyV).(*router.Value)

		//A deferred function's args are evaluated when the defer statement is evaluated
		//Deferred functions are executed in LIFO order after the surrounding function returns.
		//Deferred functions may read and assign to the returning function's named return values.
		defer func() {

			//recover is a built-in function that regains control of a panicking goroutine.
			//It is only useful inside deferred functions. In normal execution
			//it returns nil and have no other effects.
			//If the current goroutine is panicking, a call to recover
			//will capture the value given to panic and resume normal execution.
			if r := recover(); r != nil {
				log.Printf("%s: ERROR : Panic Caught : %s\n", v.RqID, r)

				//Response with error.
				router.ResponseError(ctx, w, errors.New("could not be handled"), http.StatusInternalServerError)

				//log the stack (debug.stack()) is used to format the stack trace of the
				//goroutine that calls it.
				log.Printf("%s : ERROR : Stacktrace\n", v.RqID, debug.Stack())
			}
		}()

		//Log and catch all occuring erros in middleware.
		if err := h(ctx, w, r, params); err != nil {

			err := errors.Cause(err)

			if err != router.ErrNotFound {
				log.Printf("%s ERROR : %v\n", v.RqID, err)
			}

			router.Error(ctx, w, err)

			return nil
		}
		return nil
	}

}

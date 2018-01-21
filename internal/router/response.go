package router

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

/*
DON'T JUST CHECK ERRORS HANDLE THEM GRACEFULLY. GO PROVERB.
*/

var (

	//ErrNotHealth occurs when the system doesn't work as espexted
	ErrNotHealth = errors.New("system doesn't work as espected healthy")

	//ErrNotFound is used instead of the mgo error not found
	ErrNotFound = errors.New("Entry could not be found")

	//ErrInvalidID occurs by invalid ID
	ErrInvalidID = errors.New("entered ID not valid")

	//ErrValidation occurs when invalid struct is used
	ErrValidation = errors.New("Validation error")

	//ErrDBNotConfigured when connecting to database
	ErrDBNotConfigured = errors.New("DB weither initialized or not abble to connect to")
)

//JSONError is the structure of response for errors
type JSONError struct {
	Error string `json:"error"`
	//omitempty option specifies that the field should be omitted from the encoding if the field
	//has an empty value, defined as false, 0, a nil pointer, a nil interface value, and any empty
	//array, slice, map, or string.
	Fields InvalidErrs `json:"fields, omitempty"`
}

//Error handles all errors occuring within the API
func Error(ctx context.Context, w http.ResponseWriter, err error) {

	//Cause is a method that implements interface causer. errors.Cause is used
	//to inspect all underlying errors's cause.
	switch errors.Cause(err) {
	case ErrNotFound:
		ResponseError(ctx, w, err, http.StatusNotFound)
		return

	case ErrValidation:
		ResponseError(ctx, w, err, http.StatusBadRequest)
		return

	case ErrNotHealth:
		ResponseError(ctx, w, err, http.StatusInternalServerError)
		return
	}

	//check all other kind of errors
	switch er := errors.Cause(err).(type) {
	case InvalidErrs:
		v := JSONError{
			Error:  "failed when validating struct",
			Fields: er,
		}
		Response(ctx, w, v, http.StatusBadRequest)
		return
	}
	ResponseError(ctx, w, err, http.StatusInternalServerError)
}

//Response is the json http response sent to the client after processing his request
func Response(ctx context.Context, w http.ResponseWriter, data interface{}, code int) {

	//set the status code for the request logger middleware
	//It's good to note that once we get the data of the context we need to cast it to the proper type
	//to use it. This is because the context value() function returns an interface{} type.
	v := ctx.Value(KeyV).(*Value)
	v.State = code

	//write the status code to response context
	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return
	}

	//set the content-type
	w.Header().Set("Content-type", "application/json")

	//write the status code to the response and context
	w.WriteHeader(code)

	//marshal the data into a json string.
	//MarshalIdent is like Marshal but applies Ident to format the output
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Printf("%s : Respond %v Marshalling JSON response\n", v.RqID, err)
		jsonData = []byte("{}")
	}

	//send the result back to the client
	io.WriteString(w, string(jsonData))

}

//ResponseError is used to send error response in json format.
func ResponseError(ctx context.Context, w http.ResponseWriter, err error, code int) {
	Response(ctx, w, JSONError{Error: err.Error()}, code)
}

package router

import (
	"encoding/json"
	"fmt"
	"io"

	validator "gopkg.in/go-playground/validator.v8"
)

var validate = validator.New(&validator.Config{
	TagName:      "validate",
	FieldNameTag: "json",
})

//Invalid describes a validation error belonging to a specific field.
type Invalid struct {
	Feld string `json:"field_name"`
	Err  string `json:"error"`
}

//InvalidErrs will hold all invalid fields
type InvalidErrs []Invalid

// Error implements the error interface for InvalidError.
func (err InvalidErrs) Error() string {
	var str string
	for _, v := range err {
		str = fmt.Sprintf("%s,{%s:%s}", str, v.Feld, v.Err)
	}
	return str
}

//Unmarshall decodes the input to struct type and checks the
//fileds to verify the value is valid.
func Unmarshall(r io.Reader, v interface{}) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return err
	}

	var inv InvalidErrs
	if fe := validate.Struct(v); fe != nil {
		for _, fev := range fe.(validator.ValidationErrors) {
			inv = append(inv, Invalid{Feld: fev.Field, Err: fev.Tag})
		}

		//We can return inv (type InvalidErrs) event thought this
		//method is supposed to return an error
		//because InvalidErrs implements the Error interface.
		//In design pattern in golang it calls dynamic dispatch of calls on interfaces.
		return inv
	}
	return nil
}

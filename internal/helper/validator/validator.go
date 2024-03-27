package validator

import val "github.com/go-playground/validator/v10"

var v *val.Validate

func init() {
	v = val.New()
}

func ValidateStruct(obj interface{}) error {
	return v.Struct(obj)
}

func ValidateVar(obj interface{}, tag string) error {
	return v.Var(obj, tag)
}

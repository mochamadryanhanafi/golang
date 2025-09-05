package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// DecodeAndValidate adalah helper untuk men-decode request body JSON dan memvalidasinya.
func DecodeAndValidate(r *http.Request, v interface{}, validate *validator.Validate) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return validate.Struct(v)
}

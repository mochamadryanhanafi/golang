package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	return validator.New()
}

func DecodeAndValidate(r *http.Request, v interface{}, validate *validator.Validate) error {
	log.Printf("INFO: Decoding and validating request for %s %s", r.Method, r.URL.Path)

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		log.Printf("ERROR: Failed to decode request body: %v", err)
		return fmt.Errorf("invalid request body: %w", err)
	}
	defer r.Body.Close()

	if err := validate.Struct(v); err != nil {
		var errorMsg strings.Builder
		for _, err := range err.(validator.ValidationErrors) {
			errorMsg.WriteString(fmt.Sprintf("field '%s' failed on the '%s' tag; ", err.Field(), err.Tag()))
		}
		formattedError := fmt.Errorf("validation error: %s", strings.TrimSuffix(errorMsg.String(), "; "))
		log.Printf("WARN: Validation failed for %s %s: %v", r.Method, r.URL.Path, formattedError)
		return formattedError
	}

	log.Printf("INFO: Request validation successful for %s %s", r.Method, r.URL.Path)
	return nil
}

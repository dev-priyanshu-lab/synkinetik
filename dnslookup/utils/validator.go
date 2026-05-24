package utils

import "github.com/go-playground/validator/v10"

var _validate *validator.Validate

type _fqdn struct {
	domain string `validate:"fqdn"`
}

func ValidateDomain(domain string) bool {
	fqdn := &_fqdn{domain: domain}
	if _validate == nil {
		_validate = validator.New(validator.WithRequiredStructEnabled())
	}
	err := _validate.Struct(fqdn)
	if err != nil {
		return false
	} else {
		return true
	}
}

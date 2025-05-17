package server

import (
    "errors"

    "github.com/go-playground/validator/v10"
)

// Config holds the server configuration options.
type Config struct {
    Host string `validate:"required,hostname|ip"`
    Port int    `validate:"required,min=1,max=65535"`
}

// ValidateConfig checks that required fields are present and valid.
func ValidateConfig(config *Config) error {
    v := validator.New()
    if err := v.Struct(config); err != nil {
        return errors.New("config validation failed: " + err.Error())
    }
    return nil
}
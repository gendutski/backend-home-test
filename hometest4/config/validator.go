package config

import (
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func InitValidator() *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())

	// validate username
	validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
		return re.MatchString(fl.Field().String())
	})

	return &Validator{
		validate: validate,
	}
}

type Validator struct {
	validate *validator.Validate
}

func (e *Validator) Validate(payload interface{}, errMessage string) error {
	err := e.validate.Struct(payload)
	if err != nil {
		if errMessage == "" {
			errMessage = err.Error()
		}
		return echo.NewHTTPError(http.StatusBadRequest, errMessage)
	}
	return nil
}

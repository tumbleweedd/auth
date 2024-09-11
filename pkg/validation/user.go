package validation

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// usernameRegexp represents the regular expression for a valid username.
const usernameRegexp = `^[a-zA-Z0-9_]+$`

// compiledUsernameRegexp is a compiled version of usernameRegexp for efficient matching.
var compiledUsernameRegexp = regexp.MustCompile(usernameRegexp)

// UserValidator is responsible for validating user-related data.
type UserValidator struct {
	validator *validator.Validate
}

// NewUserValidator instantiates a new UserValidator.
func NewUserValidator() *UserValidator {
	return &UserValidator{
		validator: validator.New(),
	}
}

// Validate registers a custom validation for the UserValidator.
func (u *UserValidator) Validate() error {
	if err := u.registerUsernameRegexpValidation(); err != nil {
		return fmt.Errorf("failed to register username validation: %w", err)
	}
	return nil
}

// registerUsernameRegexpValidation registers a custom validation rule for usernames.
func (u *UserValidator) registerUsernameRegexpValidation() error {
	usernameValidationFunc := func(fl validator.FieldLevel) bool {
		return compiledUsernameRegexp.MatchString(fl.Field().String())
	}

	if err := u.validator.RegisterValidation("username_regexp", usernameValidationFunc); err != nil {
		return err
	}

	return nil
}

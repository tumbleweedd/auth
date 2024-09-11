package validation

import (
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"regexp"
)

/*type Validator interface {
	Validate() error
}

type ChainValidator []Validator

func (c ChainValidator) Validate() error {
	for _, validator := range c {
		if err := validator.Validate(); err != nil {
			return err
		}
	}

	return nil
}*/

var (
	ErrInvalidType   = fmt.Errorf("invalid type")
	ErrInvalidLength = fmt.Errorf("invalid length")
	ErrEmpty         = fmt.Errorf("empty")
	ErrInvalidRange  = fmt.Errorf("invalid range")
	ErrInvalidEmail  = fmt.Errorf("invalid email")
)

type Rule func(key string, value interface{}) error

type Validator struct {
	rules map[string][]Rule
}

func New() *Validator {
	return &Validator{
		rules: make(map[string][]Rule),
	}
}

func (v *Validator) RegisterRule(key string, rule Rule) {
	v.rules[key] = append(v.rules[key], rule)
}

func (v *Validator) Validate(data interface{}) []error {
	var errors []error
	val := reflect.ValueOf(data)

	if val.Kind() != reflect.Struct {
		return []error{ErrInvalidType}
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)
		fieldName := fieldType.Name

		for _, rule := range v.rules[fieldName] {
			if err := rule(fieldName, field.Interface()); err != nil {
				errors = append(errors, err)
			}
		}
	}

	return errors
}

func ValidateLength(minLength, maxLength int) Rule {
	return func(key string, value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("%w: %s is not a string", ErrInvalidType, key)
		}

		if len(str) < minLength {
			return fmt.Errorf("%w: %s must be at least %d characters long", ErrInvalidLength, key, minLength)
		}

		if len(str) > maxLength {
			return fmt.Errorf("%w: %s must be less or equal to %d characters long", ErrInvalidLength, key, maxLength)
		}

		return nil
	}
}

func ValidatePresence() Rule {
	return func(key string, value interface{}) error {
		if value == nil || (reflect.ValueOf(value).Kind() == reflect.String && value.(string) == "") {
			return fmt.Errorf("%w: %s can't be empty", ErrEmpty, key)
		}
		return nil
	}
}

func ValidateRange(min, max int) Rule {
	return func(key string, value interface{}) error {
		num, ok := value.(int)
		if !ok {
			return fmt.Errorf("%w: %s is not a number", ErrInvalidType, key)
		}

		if num < min || num > max {
			return fmt.Errorf("%w: %s must be between %d and %d", ErrInvalidRange, key, min, max)
		}

		return nil
	}
}

func ValidateUUID() Rule {
	return func(key string, value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("%w: %s is not a string", ErrInvalidType, key)
		}

		_, err := uuid.Parse(str)
		if err != nil {
			return fmt.Errorf("%w: %s is not a valid UUID", ErrInvalidType, key)
		}

		return nil
	}
}

func ValidateEmail() Rule {
	return func(key string, value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("%w: %s is not a string", ErrInvalidType, key)
		}

		const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		emailRegex := regexp.MustCompile(emailRegexPattern)

		if !emailRegex.MatchString(str) {
			return fmt.Errorf("%w: %s is not a valid email address", ErrInvalidEmail, key)
		}

		return nil
	}
}

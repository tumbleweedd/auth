package create

import "github.com/tumbleweedd/svc/auth_service/pkg/validation"

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
}

type CreateUserResponse struct {
	UUID string `json:"uuid"`
}

func (r *CreateUserRequest) Validate() error {
	validator := validation.New()

	validator.RegisterRule("FirstName", validation.ValidatePresence())
	validator.RegisterRule("FirstName", validation.ValidateLength(2, 50))

	validator.RegisterRule("Email", validation.ValidatePresence())
	validator.RegisterRule("Email", validation.ValidateEmail())

	validator.RegisterRule("Password", validation.ValidatePresence())
	validator.RegisterRule("Password", validation.ValidateLength(8, 100))

	errors := validator.Validate(r)

	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}

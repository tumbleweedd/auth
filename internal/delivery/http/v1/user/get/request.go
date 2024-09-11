package get

import (
	"github.com/tumbleweedd/svc/auth_service/pkg/validation"
)

type GetUserRequest struct {
	UUID string `json:"uuid"`
}

func (r *GetUserRequest) Validate() error {
	validator := validation.New()

	validator.RegisterRule("UUID", validation.ValidatePresence())
	validator.RegisterRule("UUID", validation.ValidateUUID())

	errors := validator.Validate(r)

	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}

type UserResponse struct {
	UUID      string `json:"uuid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type GetUserResponse struct {
	UserResponse UserResponse `json:"user"`
}

type GetUsersRequest struct {
	UUIDs []string `json:"uuids"`
}

func (r *GetUsersRequest) Validate() error {
	validator := validation.New()

	validator.RegisterRule("UUID", validation.ValidatePresence())
	validator.RegisterRule("UUID", validation.ValidateUUID())

	errors := validator.Validate(r)

	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}

type GetUsersResponse struct {
	Users []*UserResponse `json:"users"`
}

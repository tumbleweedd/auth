package delete

import "github.com/tumbleweedd/svc/auth_service/pkg/validation"

type DeleteUserRequest struct {
	UserUUID string `json:"uuid"`
}

func (r *DeleteUserRequest) Validate() error {
	validator := validation.New()

	validator.RegisterRule("UUID", validation.ValidatePresence())
	validator.RegisterRule("UUID", validation.ValidateUUID())

	errors := validator.Validate(r)

	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}

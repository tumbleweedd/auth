package create

import (
	userEntity "github.com/tumbleweedd/svc/auth_service/internal/domain/entity/user"
)

func CreateUserRequestToDomain(req *CreateUserRequest) *userEntity.User {
	return userEntity.NewUser(
		userEntity.WithFirstName(req.FirstName),
		userEntity.WithLastName(req.LastName),
		userEntity.WithLogin(req.Login),
		userEntity.WithEmail(req.Email),
		userEntity.WithPassword(req.Password),
		userEntity.WithPhone(req.Phone),
	)
}

func NewCreateUserResponse(userUUID string) CreateUserResponse {
	return CreateUserResponse{
		UUID: userUUID,
	}
}

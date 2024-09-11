package user

import (
	userProto "github.com/tumbleweedd/auth_service_proto/gen/go/user"
	"github.com/tumbleweedd/svc/auth_service/internal/domain/usecase/user"
)

func NewCreateUserRequest(req *userProto.CreateUserRequest) user.CreateUserRequest {
	return user.NewCreateUserRequestBuilder().
		WithFirstName(req.FirstName).
		WithLastName(req.LastName).
		WithEmail(req.Email).
		WithLogin(req.Login).
		WithPassword(req.Password).
		WithPhone(req.Phone).
		Build()
}

func NewCreateUserResponse(user user.CreateUserOutput) *userProto.CreateUserResponse {
	return &userProto.CreateUserResponse{
		Id: user.ID,
	}
}

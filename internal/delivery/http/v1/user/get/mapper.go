package get

import userEntity "github.com/tumbleweedd/svc/auth_service/internal/domain/entity/user"

func GetUserRequestToDomain(request *GetUserRequest) string {
	return request.UUID
}

func NewGetUserResponse(user *userEntity.User) *GetUserResponse {
	return &GetUserResponse{
		UserResponse: UserResponse{
			UUID:      user.UUID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Login:     user.Login,
			Email:     user.Email,
			Phone:     user.Phone,
		},
	}
}

func GetUsersRequestToDomain(request *GetUsersRequest) []string {
	return request.UUIDs
}

func NewGetUsersResponse(users []*userEntity.User) *GetUsersResponse {
	userResponses := make([]*UserResponse, 0, len(users))

	for _, user := range users {
		userResponses = append(userResponses, &UserResponse{
			UUID:      user.UUID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Login:     user.Login,
			Email:     user.Email,
			Phone:     user.Phone,
		})
	}

	return &GetUsersResponse{
		Users: userResponses,
	}
}

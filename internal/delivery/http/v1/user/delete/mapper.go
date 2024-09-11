package delete

func DeleteUserRequestToDomain(req *DeleteUserRequest) string {
	return req.UserUUID
}

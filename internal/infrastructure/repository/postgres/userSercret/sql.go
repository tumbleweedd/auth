package userSercret

const (
	queryCreateUserSecret = `
	INSERT INTO user_secret (user_uuid, secret, salt, created_at)
	VALUES ($1, $2, $3, $4);
`

	queryDeleteUserSecret = `
	DELETE FROM user_secret
	WHERE user_id = $1;
	`
)

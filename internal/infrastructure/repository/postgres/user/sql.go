package user

const (
	queryCreateUser = `
	INSERT INTO user (uuid, first_name, last_name, email, phone, activated, role, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
`

	queryListUsers = `
	SELECT id, uuid, name, email, phone, activated, role, created_at, updated_at
	FROM user;
`

	queryDeleteUser = `
	DELETE FROM user
	WHERE uuid = $1;
`
	queryGetUser = `
	SELECT id, uuid, name, email, phone, activated, role, created_at, updated_at
	FROM user
	WHERE uuid = $1;
`
)

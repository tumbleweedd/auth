package userSercret

import "time"

type UserSecret struct {
	ID         string    `db:"id"`
	PasswdHash string    `db:"passwd_hash"`
	PasswdSalt string    `db:"passwd_salt"`
	CreatedAt  time.Time `db:"created_at"`
}

func NewUserSecret(id string, passwdHash string, passwdSalt string) *UserSecret {
	return &UserSecret{
		ID:         id,
		PasswdHash: passwdHash,
		PasswdSalt: passwdSalt,
	}
}

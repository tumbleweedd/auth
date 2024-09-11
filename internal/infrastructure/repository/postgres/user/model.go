package user

import (
	"database/sql"
	userEntity "github.com/tumbleweedd/svc/auth_service/internal/domain/entity/user"
	"time"
)

type User struct {
	ID        int          `db:"id"`
	UUID      string       `db:"uuid"`
	FirstName string       `db:"first_name"`
	LastName  string       `db:"last_name"`
	Email     string       `db:"email"`
	Phone     string       `db:"phone"`
	Activated bool         `db:"activated"`
	Role      string       `db:"role"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

func (u *User) ToDomain() *userEntity.User {
	var updatedAt time.Time

	if u.UpdatedAt.Valid {
		updatedAt = u.UpdatedAt.Time
	}

	return userEntity.NewUser(
		userEntity.WithID(u.ID),
		userEntity.WithUUID(u.UUID),
		userEntity.WithFirstName(u.FirstName),
		userEntity.WithLastName(u.LastName),
		userEntity.WithEmail(u.Email),
		userEntity.WithPhone(u.Phone),
		userEntity.WithActivated(u.Activated),
		userEntity.WithRole(u.Role),
		userEntity.WithCreatedAt(u.CreatedAt),
		userEntity.WithUpdatedAt(updatedAt),
	)
}

type Users []User

func (u Users) ToDomain() []*userEntity.User {
	users := make([]*userEntity.User, 0, len(u))

	for _, user := range u {
		users = append(users, user.ToDomain())
	}

	return users
}

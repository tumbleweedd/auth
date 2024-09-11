package user

import (
	"github.com/tumbleweedd/svc/auth_service/internal/domain/valueobjects"
	"time"
)

type User struct {
	ID        int
	UUID      string
	FirstName string
	LastName  string
	Email     string
	Login     string
	Password  string
	Phone     string
	Activated bool
	Role      valueobjects.Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Option func(*User)

func NewUser(options ...Option) *User {
	user := &User{
		Role:      valueobjects.UserRole,
		Activated: true,
	}

	for _, opt := range options {
		opt(user)
	}

	return user
}

func WithID(id int) Option {
	return func(u *User) {
		u.ID = id
	}
}

func WithUUID(uuid string) Option {
	return func(u *User) {
		u.UUID = uuid
	}
}

func WithFirstName(firstName string) Option {
	return func(u *User) {
		u.FirstName = firstName
	}
}

func WithLastName(lastName string) Option {
	return func(u *User) {
		u.LastName = lastName
	}
}

func WithEmail(email string) Option {
	return func(u *User) {
		u.Email = email
	}
}

func WithLogin(login string) Option {
	return func(u *User) {
		u.Login = login
	}
}

func WithPassword(password string) Option {
	return func(u *User) {
		u.Password = password
	}
}

func WithPhone(phone string) Option {
	return func(u *User) {
		u.Phone = phone
	}
}

func WithActivated(activated bool) Option {
	return func(u *User) {
		u.Activated = activated
	}
}

func WithRole(role valueobjects.Role) Option {
	return func(u *User) {
		u.Role = role
	}
}

func WithCreatedAt(createdAt time.Time) Option {
	return func(u *User) {
		u.CreatedAt = createdAt
	}
}

func WithUpdatedAt(updatedAt time.Time) Option {
	return func(u *User) {
		u.UpdatedAt = updatedAt
	}
}

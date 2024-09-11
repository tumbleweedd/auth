package token

import "time"

type Token struct {
	Token     string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

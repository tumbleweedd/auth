package token

import "time"

type Token struct {
	Token     string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type GenerateTokensOutput struct {
	AccessToken  Token
	RefreshToken Token
}

type LoginInput struct {
	Login    string
	Password string
}

type LoginOutput struct {
	AccessToken  Token
	RefreshToken Token
}

type GenerateTokensInput struct {
	Login    string
	Password string
}

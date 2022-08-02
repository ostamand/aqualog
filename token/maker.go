package token

import "time"

type CreateTokenArgs struct {
	Username string
	UserID   int64
	Duration time.Duration
}

type TokenMaker interface {
	CreateToken(args CreateTokenArgs) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}

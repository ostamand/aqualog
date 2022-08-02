package token

import (
	"errors"
	"fmt"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

var (
	ErrInvalidToken   = errors.New("token is invalid")
	ErrInvalidKeySize = errors.New("invalid key size")
)

type PasetoMaker struct {
	paseto *paseto.V2
	key    []byte
}

func NewPasetoMaker(key string) (TokenMaker, error) {
	if len(key) != chacha20poly1305.KeySize {
		err := fmt.Errorf("expected %d got %d: %w", chacha20poly1305.KeySize, len(key), ErrInvalidKeySize)
		return nil, err
	}
	maker := &PasetoMaker{
		paseto: paseto.NewV2(),
		key:    []byte(key),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(args CreateTokenArgs) (string, *Payload, error) {
	payload, err := NewPayload(args.Username, args.UserID, args.Duration)
	if err != nil {
		return "", payload, err
	}
	token, err := maker.paseto.Encrypt(maker.key, payload, nil)
	return token, payload, err
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.key, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}

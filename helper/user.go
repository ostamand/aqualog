package helper

import (
	"context"

	db "github.com/ostamand/aqualog/db/sqlc"
	"github.com/ostamand/aqualog/storage"
	"github.com/ostamand/aqualog/util"
)

var PasswordMinLength = 6

type SaveUserParams struct {
	Username string
	Email    string
	Password string
}

func SaveUser(ctx context.Context, s storage.Storage, args SaveUserParams) (user db.User, err error) {
	if len(args.Password) < PasswordMinLength {
		return user, ErrPasword{Info: PasswordLength}
	}
	hashPassword, err := util.HashPassword(args.Password)
	if err != nil {
		return
	}
	user, err = s.CreateUser(ctx, db.CreateUserParams{
		Username:       args.Username,
		Email:          args.Email,
		HashedPassword: hashPassword,
	})
	return
}

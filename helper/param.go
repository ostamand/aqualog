package helper

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/ostamand/aqualog/db/sqlc"
	"github.com/ostamand/aqualog/storage"
)

type SaveParamArgs struct {
	UserID    int64
	ParamName string
	Value     float64
}

func SaveParam(ctx context.Context, store storage.Storage, args SaveParamArgs) (db.Param, error) {
	var param db.Param
	var paramType db.ParamType
	var err error

	// get param type by name
	paramType, err = store.GetParamTypeByName(ctx, db.GetParamTypeByNameParams{
		UserID: args.UserID,
		Name:   args.ParamName,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			paramType, err = store.CreateValueType(ctx, db.CreateValueTypeParams{
				Name:   args.ParamName,
				UserID: args.UserID,
			})
			if err != nil {
				return param, err
			}
		} else {
			return param, err
		}
	}

	param, err = store.CreateParam(ctx, db.CreateParamParams{
		UserID:      args.UserID,
		ParamTypeID: paramType.ID,
		Value:       args.Value,
	})

	return param, err
}

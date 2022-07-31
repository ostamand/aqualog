package helper

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/ostamand/aqualog/db/sqlc"
	"github.com/ostamand/aqualog/storage"
)

type SaveParamArgs struct {
	userID    int64
	paramName string
	value     float64
}

func SaveParam(ctx context.Context, s storage.Storage, args SaveParamArgs) (db.Param, error) {
	var param db.Param
	var paramType db.ParamType
	var err error

	// get param type by name
	paramType, err = s.GetParamTypeByName(ctx, db.GetParamTypeByNameParams{
		UserID: args.userID,
		Name:   args.paramName,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			paramType, err = s.CreateValueType(ctx, db.CreateValueTypeParams{
				Name:   args.paramName,
				UserID: args.userID,
			})
			if err != nil {
				return param, err
			}
		} else {
			return param, err
		}
	}

	param, err = s.CreateParam(ctx, db.CreateParamParams{
		UserID:      args.userID,
		ParamTypeID: paramType.ID,
		Value:       args.value,
	})

	return param, err
}

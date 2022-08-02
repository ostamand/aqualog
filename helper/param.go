package helper

import (
	"context"
	"database/sql"
	"errors"
	"time"

	db "github.com/ostamand/aqualog/db/sqlc"
	"github.com/ostamand/aqualog/storage"
)

type SaveParamArgs struct {
	UserID    int64
	ParamName string
	Value     float64
	Timestamp time.Time
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
		if errors.Is(err, sql.ErrNoRows) { // TODO this is specific to SQL interface implementation
			paramType, err = store.CreateParamType(ctx, db.CreateParamTypeParams{
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

	// set timestamp if not defined
	if args.Timestamp.IsZero() {
		args.Timestamp = time.Now()
	}

	param, err = store.CreateParam(ctx, db.CreateParamParams{
		UserID:      args.UserID,
		ParamTypeID: paramType.ID,
		Value:       args.Value,
		Timestamp:   args.Timestamp,
	})

	return param, err
}

type GetParamsArgs struct {
	From          time.Time
	To            time.Time
	UserID        int64
	ParamTypeName string
	Limit         int32
	Offset        int32
}

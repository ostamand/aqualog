// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateParam(ctx context.Context, arg CreateParamParams) (Param, error)
	CreateParamType(ctx context.Context, arg CreateParamTypeParams) (ParamType, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteParam(ctx context.Context, arg DeleteParamParams) error
	DeleteUser(ctx context.Context, id int64) error
	GetParam(ctx context.Context, id int64) (Param, error)
	GetParamType(ctx context.Context, id int64) (ParamType, error)
	GetParamTypeByName(ctx context.Context, arg GetParamTypeByNameParams) (ParamType, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetUser(ctx context.Context, username string) (User, error)
	ListParamOrigins(ctx context.Context, paramTypeName string) ([]ParamOrigin, error)
	ListParamsByType(ctx context.Context, arg ListParamsByTypeParams) ([]ListParamsByTypeRow, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateParamType(ctx context.Context, arg UpdateParamTypeParams) (ParamType, error)
}

var _ Querier = (*Queries)(nil)

package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/ostamand/aqualog/db/sqlc"
	"github.com/ostamand/aqualog/token"
)

type updateParamTypeRequest struct {
	Target float64 `json:"target"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
}

func (server *Server) updateParamType(ctx *gin.Context) {
	var req updateParamTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	payload, ok := ctx.MustGet(authPayloadKey).(*token.Payload)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, token.ErrInvalidToken)
		return
	}
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// TODO does not manage setting to null
	paramType, err := server.storage.UpdateParamType(ctx, db.UpdateParamTypeParams{
		UserID: payload.UserID,
		ID:     id,
		Target: sql.NullFloat64{
			Float64: req.Target,
			Valid:   true,
		},
		Min: sql.NullFloat64{
			Float64: req.Min,
			Valid:   true,
		},
		Max: sql.NullFloat64{
			Float64: req.Max,
			Valid:   true,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, paramType)
}

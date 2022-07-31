package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ostamand/aqualog/helper"
	"github.com/ostamand/aqualog/token"
)

type createParamRequest struct {
	Value     float64 `json:"value" binding:"required,min=0"`
	ParamType string  `json:"type" binding:"required"`
}

func (server *Server) createParam(ctx *gin.Context) {
	var req createParamRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload, ok := ctx.MustGet(authPayloadKey).(*token.Payload)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, token.ErrInvalidToken)
		return
	}
	param, err := helper.SaveParam(ctx, server.storage, helper.SaveParamArgs{
		UserID:    authPayload.UserID,
		ParamName: req.ParamType,
		Value:     req.Value,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, param)
}

type getValueRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getValue(ctx *gin.Context) {
	var req getValueRequest
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	value, err := server.storage.GetParam(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, value)
}

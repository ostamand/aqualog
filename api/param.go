package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/ostamand/aqualog/db/sqlc"
	"github.com/ostamand/aqualog/helper"
	"github.com/ostamand/aqualog/token"
)

type createParamRequest struct {
	Value     float64   `json:"value" binding:"required,min=0"`
	ParamType string    `json:"type" binding:"required"`
	Timestamp time.Time `json:"timestamp"`
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
		Timestamp: req.Timestamp,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, param)
}

type getParamsRequest struct {
	ParamType string    `form:"type"`
	Limit     int32     `form:"limit,default=100"`
	Offset    int32     `form:"offset,default=0"`
	From      time.Time `form:"from"`
	To        time.Time `form:"to"`
}

func (server *Server) getParams(ctx *gin.Context) {
	var req getParamsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	payload, ok := ctx.MustGet(authPayloadKey).(*token.Payload)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, token.ErrInvalidToken)
		return
	}
	args := db.ListParamsByTypeParams{
		UserID:        payload.UserID,
		ParamTypeName: req.ParamType, // TODO what to do when not provided?
		Limit:         req.Limit,
		Offset:        req.Offset,
		From:          req.From,
		To:            req.To,
	}
	args.FillDefaults()
	params, err := server.storage.ListParamsByType(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, params)
}

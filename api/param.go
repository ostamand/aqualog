package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createValueRequest struct {
	Value     float64 `json:"value" binding:"required,min=0"`
	ValueType string  `json:"type" binding:"required"`
}

func (server *Server) createValue(ctx *gin.Context) {
	var req createValueRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	//authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)

	//db.CreateValueParams{}
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

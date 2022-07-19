package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ostamand/aqualog/storage"
)

type saveValueRequest struct {
	Username string  `json:"username" binding:"required"`
	Value    float64 `json:"value" binding:"required"`
	Type     string  `json:"type" binding:"required"`
}

func (server *Server) saveValue(ctx *gin.Context) {
	var req saveValueRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := storage.AddValueParams{
		Username: req.Username,
		Value:    req.Value,
		Type:     req.Type,
	}
	value, err := server.storage.AddValue(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, value)
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
	value, err := server.storage.GetValue(ctx, req.ID)
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

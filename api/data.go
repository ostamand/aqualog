package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getDataRequest struct {
	DataType      string `form:"type"`
	ParamTypeName string `form:"param"`
}

func (server *Server) getData(ctx *gin.Context) {
	var req getDataRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	// for now, no need for an auth token. might change later
	switch req.DataType {
	case "origins":
		origins, err := server.storage.ListParamOrigins(ctx, req.ParamTypeName)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, origins)
		return
	default:
		err := fmt.Errorf("type not supported")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}

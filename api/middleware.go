package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ostamand/aqualog/token"
)

const (
	authHeaderKey  = "authorization"
	authType       = "bearer"
	authPayloadKey = "auth_payload"
)

var errInvalidToken = errorResponse(fmt.Errorf("invalid authorization header"))

func corsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
		ctx.Next()
	}
}

func authMiddleware(token token.TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get token
		authHeader := ctx.GetHeader(authHeaderKey)

		// check token format
		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errInvalidToken)
			return
		}
		if strings.ToLower(fields[0]) != authType {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errInvalidToken)
			return
		}
		accessToken := fields[1]

		// verify token
		payload, err := token.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		}

		// set user payload
		ctx.Set(authPayloadKey, payload)
		ctx.Next()
	}
}

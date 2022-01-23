package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/harrychopra/go-api/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenmaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}
		authorizationnType := strings.ToLower(fields[0])
		if authorizationnType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}
		accessToken := fields[1]
		payload, err := tokenmaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next() // Forward the request to next handler
	}
}

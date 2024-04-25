package api

import (
	"errors"
	"fmt"
	"net/http"
	"simplebank/token"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer" //假设银行现在只支持一种不记名验证方式
	authorizationPayLoadKey = "authorization_payLoadKey"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader("authorizationHeaderKey")
		fmt.Println(authorizationHeader, "1111111")
		if len(authorizationHeader) == 0 {
			err := errors.New("无效的验证")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		fields := strings.Fields(authorizationHeader) //获取验证方式和验证token
		if len(fields) < 2 {
			err := errors.New("不合法的格式")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer { //假设不支持这种验证方式
			err := fmt.Errorf("不支持这种授权类型%s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		acessToken := fields[1]
		payLoad, err := tokenMaker.VerifyToken(acessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.Set(authorizationPayLoadKey, payLoad)
		ctx.Next()

	}

}

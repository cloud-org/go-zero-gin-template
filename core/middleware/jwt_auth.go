package middleware

import (
	"errors"
	"net/http"

	"go-zero-gin-template/core/common"
	"go-zero-gin-template/core/jwt"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/gin-gonic/gin"
)

func AssertUserInfo(ctx *gin.Context) (*jwt.CustomClaims, error) {
	value, ok := ctx.Get("user")
	if !ok {
		return nil, errors.New("context 不存在 user 键值")
	}

	switch user := value.(type) {
	case *jwt.CustomClaims:
		return user, nil
	default:
		return nil, errors.New("断言 ctx.user 失败")
	}
}

func JWTAuth(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 通过http header中的token解析来认证
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusOK, common.RespWithError("no token"))
			c.Abort()
			return
		}

		// 初始化一个JWT对象,并根据结构体方法来解析token
		j := jwt.NewJWT(key)

		// 解析token中包含的相关信息
		if claims, err := j.ParserToken(token); err != nil {
			err1 := j.ExpiredTokenError()
			if errors.As(err, &err1) {
				c.JSON(http.StatusOK, common.RespWithErrorCode(common.ErrCodeTokenAuthErr, "token 已过期"))
				c.Abort()
				return
			}
			c.JSON(http.StatusOK, common.RespWithErrorCode(common.ErrCodeTokenExpired, "token 校验失败"))
			c.Abort()
			return
		} else {
			// 将解析后的有效载荷claims重新写入gin.context引用对象中
			logx.Infof("user: %v, expire at: %v", claims.UserName, claims.ExpiresAt)
			c.Set("user", claims)
			c.Next()
		}
	}
}

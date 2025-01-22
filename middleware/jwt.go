package middleware

import (
	"Fire/pkg/e"
	"Fire/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = 200

		token := c.GetHeader("Authorization")
		if token == "" {
			code = e.ErrorAuth
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				if time.Now().Unix() > claims.ExpiresAt {
					code = e.ErrorAuthCheckTokenTimeout
				} else {
					code = e.ErrorAuthCheckTokenFail
				}
			}
		}
		if code != e.SUCCESS {
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"data": nil,
				"msg":  e.GetMsg(code),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

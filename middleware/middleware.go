package middleware

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "*")
		ctx.Header("Access-Control-Allow-Headers", "*")
		ctx.Header("Access-Control-Expose-Headers", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

func Auth(force bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		flag := false
		Auth := ctx.GetHeader("Authorization")
		token := strings.TrimPrefix(Auth, "Bearer ")
		if token == "" {
			flag = true
		}

		var jwtToken *jwt.Token
		var err error
		if !flag {
			jwtToken, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				byteSlice, _ := hex.DecodeString("0x75001ea42539b2dd087f53eec22a714b4fc7cfdfd2c408914315b2ba20c05108a3b67bac62d5fc2ddf4db7f2094a6be50375e8d82abab650746ad4ddd1e1963c")
				return byteSlice, nil
			})
			if err != nil {
				flag = true
			}
		}

		if !flag {
			if _, ok := jwtToken.Claims.(jwt.MapClaims); !(ok && jwtToken.Valid) {
				flag = true
			}
		}

		if !flag {
			userId := int64(jwtToken.Claims.(jwt.MapClaims)["userId"].(float64))
			ctx.Set("userId", userId)
		} else {
			if force {
				ctx.JSON(http.StatusUnauthorized, gin.H{"code": "200", "msg": "Unauthorized"})
				ctx.Abort()
				return
			}
		}
		ctx.Next()
	}
}

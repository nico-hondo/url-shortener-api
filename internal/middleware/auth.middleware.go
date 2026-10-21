package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nico-hondo/url-shortener/internal/dto"
	"github.com/nico-hondo/url-shortener/pkg"
	"github.com/redis/go-redis/v9"
)

func CheckToken(c *gin.Context) {
	bearer := c.GetHeader("Authorization")
	if bearer == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Msg:     "please login first",
		})
		c.Abort()
		return
	}
	result := strings.Split(bearer, " ")
	if len(result) != 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Msg:     "invalid bearer token",
		})
		c.Abort()
		return
	}
	if result[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Msg:     "invalid bearer token",
		})

		return
	}
	var token pkg.JwtClaims
	err := token.DecodeToken(result[1])
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenInvalidIssuer) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Msg:     "invalid token",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Msg:     "terjadi kesalahan sistem",
		})
		return
	}
	c.Set("token", token)
	c.Next()
}

func AuthMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Msg:     "please login first",
			})
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Msg:     "invalid bearer token format",
			})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// if tokenString == authHeader {
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
		// 		Success: false,
		// 		Msg:     "invalid bearer token format",
		// 	})
		// 	return
		// }

		// 1. Decode & Validasi Token JWT
		var tokenClaims pkg.JwtClaims
		err := tokenClaims.DecodeToken(tokenString)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenInvalidIssuer) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
					Success: false,
					Msg:     "invalid or expired token",
				})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
				Success: false,
				Msg:     "terjadi kesalahan sistem",
			})
			return
		}

		// 2. Cek Blacklist di Redis berdasarkan JTI
		key := fmt.Sprintf("%s:bl:%s", os.Getenv("RDB_USER"), tokenClaims.ID)
		exists, err := rdb.Exists(ctx.Request.Context(), key).Result()
		if err != nil {
			log.Println("Redis Error Middleware: ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
				Success: false,
				Msg:     "terjadi kesalahan pada sistem otentikasi",
			})
			return
		}

		if exists > 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Msg:     "token has been revoked/logged out. Please login again.",
			})
			return
		}

		// 3. Set data claims ke context Gin
		ctx.Set("token", tokenClaims)
		ctx.Next()
	}
}

package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nico-hondo/url-shortener/internal/dto"
	"github.com/nico-hondo/url-shortener/internal/middleware"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitMainRouter(router *gin.Engine, db *pgxpool.Pool, rds *redis.Client) {
	// Pasang middleware CORS di sini secara global
	router.Use(middleware.Cors)

	router.GET("/documentation/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	initAuthRouter(router, db, rds)
	initLinkRouter(router, db)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Msg:     "Page not found!",
		})
	})
}

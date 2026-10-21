package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nico-hondo/url-shortener/internal/controller"
	"github.com/nico-hondo/url-shortener/internal/middleware"
	"github.com/nico-hondo/url-shortener/internal/repo"
	"github.com/nico-hondo/url-shortener/internal/service"
	"github.com/redis/go-redis/v9"
)

func initAuthRouter(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authRouter := r.Group("/auth")

	ar := repo.NewAuthRepo(db, rdb)
	as := service.NewAuthService(ar)
	ah := controller.NewAuthHandler(as)

	authRouter.POST("/register", ah.Register)
	authRouter.POST("/login", ah.Login)

	authRouter.POST("/logout", middleware.AuthMiddleware(rdb), ah.Logout)
	// authRouter.PATCH("/change-password", middleware.CheckToken, ah.ChangePassword)
	// authRouter.PATCH("/change-profile", middleware.CheckToken, ah.ChangeProfileUser)

	// protect := r.Group("/")
	// protect.Use(middleware.AuthMiddleware(redisClient, "YOUR_JWT_SECRET"))
	// {
	// 	protect.POST("/logout", ah.LogoutHandler)
	// }
}

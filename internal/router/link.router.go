package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nico-hondo/url-shortener/internal/controller"
	"github.com/nico-hondo/url-shortener/internal/middleware"
	"github.com/nico-hondo/url-shortener/internal/repo"
	"github.com/nico-hondo/url-shortener/internal/service"
)

func initLinkRouter(r *gin.Engine, db *pgxpool.Pool) {
	linkRouter := r.Group("/api")

	lr := repo.NewLinkRepo(db)
	ls := service.NewLinkService(lr)
	lh := controller.NewLinkController(ls)

	linkRouter.POST("/links", middleware.CheckToken, lh.CreateLinkController)
	linkRouter.GET("/links", middleware.CheckToken, lh.GetUserLinksController)
	linkRouter.DELETE("/links/:id", middleware.CheckToken, lh.DeleteLinkController)
	r.GET("/:slug", lh.RedirectLinkController)
}

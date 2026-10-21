package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nico-hondo/url-shortener/internal/dto"
	"github.com/nico-hondo/url-shortener/internal/service"
	"github.com/nico-hondo/url-shortener/pkg"
)

type LinkController struct {
	ls *service.LinkService
}

func NewLinkController(ls *service.LinkService) *LinkController {
	return &LinkController{
		ls: ls,
	}
}

// CreateLinkController godoc
// @Summary      Create Short Link
// @Description  Generate short URL with optional custom slug
// @Tags         Link
// @Accept       json
// @Produce      json
// @Security     BearerToken
// @Param        body  body      dto.CreateLinkRequest  true  "Create link body"
// @Success      201   {object}  dto.Response
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /api/links [post]
func (lh *LinkController) CreateLinkController(ctx *gin.Context) {
	var body dto.CreateLinkRequest

	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Msg:     "not valid input data: " + err.Error(),
		})
		return
	}

	val, exists := ctx.Get("token")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Msg:     "unauthorized",
		})
		return
	}

	claims, ok := val.(pkg.JwtClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Msg:     "terjadi kesalahan sistem",
		})
		return
	}

	res, err := lh.ls.CreateLink(ctx.Request.Context(), claims.Id, body)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Msg:     "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Msg:     "Link was Shortened",
		Data:    res,
	})
}

// RedirectLinkController godoc
// @Summary      Redirect Short Link
// @Description  Redirect to original URL by slug
// @Tags         Link
// @Param        slug  path  string  true  "Short link slug"
// @Success      302
// @Failure      404   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /{slug} [get]
func (lh *LinkController) RedirectLinkController(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Msg:     "slug is required",
		})
		return
	}

	originalURL, err := lh.ls.GetOriginalURL(ctx.Request.Context(), slug)
	if err != nil {
		log.Println("Redirect error:", err.Error())
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Msg:     "URL not found",
		})
		return
	}

	// HTTP Status 302 Found untuk Temporary Redirect
	ctx.Redirect(http.StatusFound, originalURL)
}

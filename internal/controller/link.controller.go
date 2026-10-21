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

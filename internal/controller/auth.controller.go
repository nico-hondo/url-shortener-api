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

type AuthHandler struct {
	as *service.AuthService
}

func NewAuthHandler(as *service.AuthService) *AuthHandler {
	return &AuthHandler{
		as: as,
	}
}

// SignInAccount godoc
// @Summary		Login an Account
// @Description	Login Acc from body
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Param		body	body		dto.LoginRequest	true	"body to login"
// @Success		200		{object}	dto.Response
// @Failure		400		{object}	dto.ErrorResponse
// @Failure		401		{object}	dto.ErrorResponse
// @Failure		500		{object}	dto.ErrorResponse
// @Router		/auth/login	[post]
func (a *AuthHandler) Login(c *gin.Context) {
	var body dto.LoginRequest
	if err := c.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Msg:     "internal server error",
		})
		return
	}

	// gunakan service
	token, user, err := a.as.LoginService(c.Request.Context(), body)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data: gin.H{
			"token": token,
			"user":  user,
		},
		Msg: "Login Success",
	})
}

// RegistrationAccount godoc
// @Summary      Register an Account
// @Description  Regist Acc from body
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RegisterRequest  true  "body to register"
// @Success      201   {object}  dto.Response
// @Failure 	 401   {object}  dto.Response
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /auth/register [post]
func (ah *AuthHandler) Register(ctx *gin.Context) {
	var body dto.RegisterRequest

	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err.Error())

		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Msg:     "not valid input data: " + err.Error(),
		})
		return
	}

	// gunakan service
	if err := ah.as.RegisterService(ctx.Request.Context(), body); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Msg:     "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Msg:     "Register Success",
	})
}

// Logout 		godoc
// @Summary 	Log out an Account
// @Description Blacklist the authenticated access token in Redis until its expiration.
// @Tags 		Auth
// @Accept 		json
// @Produce 	json
// @Security 	BearerToken
// @Success 	200 {object} dto.Response "Successfully logged out"
// @Failure 	401 {object} dto.Response "Not authenticated"
// @Failure 	500 {object} dto.Response "Server error"
// @Router 		/auth/logout [post]
func (ar *AuthHandler) Logout(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Msg:     "Token Data Not Found in Context.",
		})
		return // Penting: cegah execution lanjut ke bawah biar gak panic
	}

	claims, ok := token.(pkg.JwtClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Msg:     "Invalid Token Claims Format",
		})
		return
	}

	jti := claims.ID

	if claims.ExpiresAt == nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Msg:     "Token does not have expiration claim",
		})
		return
	}
	expiresAt := claims.ExpiresAt.Time

	if err := ar.as.Logout(ctx.Request.Context(), jti, expiresAt); err != nil {
		log.Println("Error logging out: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Msg:     "Successfully logged out",
	})
}

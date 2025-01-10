package controller

import (
	"net/http"
	"scs-session/internal/dto"
	"scs-session/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthController(a usecase.AuthUsecase) AuthController {
	return AuthController{
		authUsecase: a,
	}
}

func (a *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		response := dto.BaseResponse{Message: "error parse"}
		c.Status(400)
		c.Set("data", response)
		return
	}
	data, err := a.authUsecase.Login(c.Request.Context(), req)
	if err != nil {
		response := dto.BaseResponse{Message: err.Error()}
		c.Status(401)
		c.Set("data", response)
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    data.BlankToken,
		MaxAge:   int(time.Until(data.BlankTokenExpired).Seconds()),
		Domain:   "localhost",
		Path:     "/",
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	})
	c.Set("data", data)
}

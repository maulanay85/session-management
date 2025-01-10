package controller

import (
	"scs-session/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUsecase usecase.UserUsecase
}

func NewUserController(a usecase.UserUsecase) UserController {
	return UserController{
		userUsecase: a,
	}
}

func (a *UserController) GetProfile(c *gin.Context) {
	data, err := a.userUsecase.FindById(c.Request.Context())
	if err != nil {
		c.Status(400)
		return
	}
	c.Status(200)
	c.Set("data", data)
}

package controller

import (
	"scs-session/internal/dto"
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

func (a *UserController) Update(c *gin.Context) {
	id := c.Param("id")
	var request dto.UserUpdateRequest
	if err := c.BindJSON(&request); err != nil {
		response := dto.BaseResponse{Message: "error parse"}
		c.Status(400)
		c.Set("data", response)
		return
	}
	err := a.userUsecase.Update(c.Request.Context(), id, request)
	if err != nil {
		response := dto.BaseResponse{Message: err.Error()}
		c.Status(500)
		c.Set("data", response)
		return
	}
	c.Status(201)
}

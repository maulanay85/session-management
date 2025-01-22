package controller

import (
	"scs-session/internal/domain"
	"scs-session/internal/dto"
	"scs-session/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AuditTrailController struct {
	au usecase.AuditTrailUsecase
}

func NewAuditTrailController(au usecase.AuditTrailUsecase) AuditTrailController {
	return AuditTrailController{
		au: au,
	}
}

func (a *AuditTrailController) GetAll(c *gin.Context) {
	userId := c.Query("userId")
	serviceName := c.Query("serviceName")

	data, err := a.au.GetAll(c.Request.Context(), domain.AuditTrailFilter{
		EntityID:    userId,
		ServiceName: serviceName,
	})
	if err != nil {
		c.Status(500)
		c.Set("data", dto.BaseResponse{Message: err.Error()})
		return
	}
	c.Status(200)
	c.Set("data", dto.BaseResponse{Data: data})
}

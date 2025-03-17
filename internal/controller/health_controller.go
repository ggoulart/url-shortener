package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthService interface {
	Health(ctx context.Context) map[string]bool
}

type HealthController struct {
	service HealthService
}

func NewHealthController(service HealthService) *HealthController {
	return &HealthController{service: service}
}

func (h *HealthController) Health(ctx *gin.Context) {
	health := h.service.Health(ctx)

	ctx.JSON(http.StatusOK, health)
}

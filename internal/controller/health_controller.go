package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthService interface {
	Health(ctx context.Context) (map[string]bool, error)
}

type HealthController struct {
	service HealthService
}

func NewHealthController(service HealthService) *HealthController {
	return &HealthController{service: service}
}

func (h *HealthController) Health(ctx *gin.Context) {
	health, err := h.service.Health(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, health)
}

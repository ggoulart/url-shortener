package middleware

import (
	"errors"
	"net/http"

	"github.com/ggoulart/url-shortener/internal/controller"
	"github.com/ggoulart/url-shortener/internal/repository"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if err := c.Errors.Last(); err != nil {
			var status int

			switch {
			case errors.Is(err.Err, controller.ErrBadRequest),
				errors.Is(err.Err, repository.ErrNotFound):
				status = http.StatusBadRequest
			default:
				status = http.StatusInternalServerError
			}

			c.JSON(status, gin.H{"error": err.Error()})
		}
	}
}

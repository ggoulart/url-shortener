package controller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

var ErrBadRequest = errors.New("invalid body")

type ShortenerService interface {
	Shortener(ctx context.Context, longURL string) (url.URL, error)
	Retrieve(ctx context.Context, encodedKey string) (url.URL, error)
}

type ShortenerController struct {
	service ShortenerService
}

func NewShortenerController(service ShortenerService) *ShortenerController {
	return &ShortenerController{service: service}
}

func (c *ShortenerController) ShortenURL(ctx *gin.Context) {
	var body ShortenerRequest
	err := ctx.BindJSON(&body)
	if err != nil {
		slog.Warn(fmt.Sprintf("failed to parse body: %v", err))
		ctx.Error(ErrBadRequest)
		return
	}

	shortURL, err := c.service.Shortener(ctx, body.LongURL)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, ShortenerResponse{ShortURL: shortURL.String()})
}

func (c *ShortenerController) RetrieveURL(ctx *gin.Context) {
	encodedKey := ctx.Param("encodedKey")

	longURL, err := c.service.Retrieve(ctx, encodedKey)
	if err != nil {
		ctx.Error(err)
		return
	}

	http.Redirect(ctx.Writer, ctx.Request, longURL.String(), http.StatusFound)
}

type ShortenerRequest struct {
	LongURL string `json:"longUrl" binding:"required"`
}

type ShortenerResponse struct {
	ShortURL string `json:"shortUrl" binding:"required"`
}

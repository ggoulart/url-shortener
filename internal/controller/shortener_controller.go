package controller

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

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
		e := fmt.Errorf("failed parse request body: %v", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return
	}

	shortURL, err := c.service.Shortener(ctx, body.LongURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, ShortenerResponse{ShortURL: shortURL.String()})
}

func (c *ShortenerController) RetrieveURL(ctx *gin.Context) {
	encodedKey := ctx.Param("encodedKey")

	longURL, err := c.service.Retrieve(ctx, encodedKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	http.Redirect(ctx.Writer, ctx.Request, longURL.String(), http.StatusFound)
}

type ShortenerRequest struct {
	LongURL string `json:"longUrl" binding:"required"`
}

type ShortenerResponse struct {
	ShortURL string `json:"shortUrl" binding:"required"`
}

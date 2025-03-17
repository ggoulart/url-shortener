package main

import (
	"fmt"
	"log"

	"github.com/ggoulart/url-shortener/internal/clients/postgres"
	"github.com/ggoulart/url-shortener/internal/controller"
	"github.com/ggoulart/url-shortener/internal/repository"
	"github.com/ggoulart/url-shortener/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func main() {
	err, shortenerHost := loadConfigs()

	postgresConfig, err := postgres.NewConfig()
	if err != nil {
		log.Panic(err)
	}

	postgresClient, err := postgres.NewClient(*postgresConfig)
	if err != nil {
		log.Panic(err)
	}
	defer postgresClient.DB.Close()

	shortenerRepository := repository.NewShortenerRepository(postgresClient.DB)
	shortenerService := service.NewShortenerService(shortenerRepository, shortenerHost, uuid.New().String)
	shortenerController := controller.NewShortenerController(shortenerService)

	healthService := service.NewHealthService(postgresClient)
	healthController := controller.NewHealthController(healthService)

	r := gin.Default()

	routes(r, shortenerController, healthController)

	err = r.Run(":8080")
	if err != nil {
		log.Panic(fmt.Errorf("failed to start server: %v", err))
	}
}

func loadConfigs() (error, string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(fmt.Errorf("failed to load config file: %s", err))
	}

	shortenerHost := viper.GetString("service.SHORTENER_HOST")
	return err, shortenerHost
}

func routes(r *gin.Engine, shortenerController *controller.ShortenerController, healthController *controller.HealthController) {
	r.POST("/api/v1/shorten", shortenerController.ShortenURL)
	r.GET("/api/v1/:encodedKey", shortenerController.RetrieveURL)
	r.GET("/api/v1/health", healthController.Health)
}

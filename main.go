package main

import (
	"net/http"

	"github.com/RedHatInsights/quickstarts/config"
	"github.com/RedHatInsights/quickstarts/pkg/database"
	"github.com/RedHatInsights/quickstarts/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func initDependecies() {
	database.Init()
}

func main() {
	godotenv.Load()
	config.Init()
	initDependecies()
	cfg := config.Get()
	logrus.WithFields(logrus.Fields{
		"ServerAddr": cfg.ServerAddr,
	})

	// done := make(chan struct{})
	// sigint := make(chan os.Signal, 1)
	// signal.Notify(sigint)

	engine := gin.Default()
	engine.GET("/test", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "This is a test response",
		})
	})

	engine.GET("/api/quickstarts/v1/openapi.json", func(c *gin.Context) {
		c.File(cfg.OpenApiSpecPath)
	})

	engine.POST("/api/quickstarts/v1/quickstarts", func(c *gin.Context) {
		var quickStart *models.Quickstart
		if err := c.ShouldBindJSON(&quickStart); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		}

		database.DB.Create(&quickStart)
		c.JSON(http.StatusOK, gin.H{"id": quickStart.ID})
	})

	engine.GET("/api/quickstarts/v1/quickstarts", func(c *gin.Context) {
		var quickStarts []models.Quickstart
		database.DB.Find(&quickStarts)
		c.JSON(http.StatusOK, gin.H{"data": quickStarts})
	})

	server := http.Server{
		Addr:    cfg.ServerAddr,
		Handler: engine,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatal("listen: %s\n", err)
	}

	// <-done
	// logrus.Info("Gracefully stopping server")

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// defer func() {
	// 	// extra handling here
	// 	cancel()
	// }()

	// if err := server.Shutdown(ctx); err != nil {
	// 	logrus.Fatal("Server shutdown failed:%+v", err)
	// }
	// logrus.Info("Server stypped properly")
}

package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"project-app-ecommerce-golang-tim-1/pkg/utils"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ApiServer(config utils.Configuration, logger *zap.Logger, h *gin.Engine) {
	// register order and address routes (handlers to be wired in wiring layer)
	h.POST("/order", func(c *gin.Context) { c.JSON(501, gin.H{"message": "not implemented"}) })
	h.GET("/cart", func(c *gin.Context) { c.JSON(501, gin.H{"message": "not implemented"}) })
	// address CRUD
	h.POST("/address", func(c *gin.Context) { c.JSON(501, gin.H{"message": "not implemented"}) })
	h.GET("/address", func(c *gin.Context) { c.JSON(501, gin.H{"message": "not implemented"}) })
	h.PUT("/address/:id", func(c *gin.Context) { c.JSON(501, gin.H{"message": "not implemented"}) })
	h.DELETE("/address/:id", func(c *gin.Context) { c.JSON(501, gin.H{"message": "not implemented"}) })
	fmt.Println("Server running on port 8080")

	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: h,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("can't run service", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server Shutdown:", zap.Error(err))
	}
	logger.Info("Server exiting")
}

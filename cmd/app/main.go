package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KAA295/lo/api/rest"
	ramstorage "github.com/KAA295/lo/internal/repository/ram_storage"
	"github.com/KAA295/lo/internal/service"
	"github.com/KAA295/lo/pkg"
)

func main() {
	logger := pkg.NewLogger(100)
	logger.Start()
	defer logger.Stop()

	taskRepo := ramstorage.NewTaskRepo()
	taskService := service.NewTaskService(taskRepo)
	taskHandler := rest.NewTaskHandler(taskService, logger)

	mu := http.NewServeMux()
	mu = taskHandler.RegisterRoutes(mu)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mu,
	}

	go func() {
		logger.Log("INFO", "Server starting")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log("ERROR", fmt.Sprintf("Server failed: %v", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log("INFO", "Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Log("ERROR", fmt.Sprintf("Server shutdown failed: %v", err))
	}

	logger.Log("INFO", "Server stopped")
}

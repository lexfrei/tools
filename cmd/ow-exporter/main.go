package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	defaultPort = "9420"
)

func main() {
	// Setup structured logging
	programLevel := new(slog.LevelVar)
	programLevel.Set(slog.LevelInfo)
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})))

	slog.Info("ow-exporter starting", "version", "development")

	// Create Echo instance
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	setupRoutes(e)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Start server
	go func() {
		slog.Info("starting HTTP server", "port", port)
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	slog.Info("shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server shutdown complete")
}

func setupRoutes(e *echo.Echo) {
	// Health check
	e.GET("/health", healthHandler)

	// API routes
	api := e.Group("/api")
	api.GET("/users", listUsersHandler)
	api.POST("/users", createUserHandler)
	api.GET("/users/:username", getUserHandler)
	api.PUT("/users/:username", updateUserHandler)
	api.DELETE("/users/:username", deleteUserHandler)

	// Prometheus metrics
	e.GET("/metrics", metricsHandler)

	// Development info
	e.GET("/", indexHandler)
}

// Handlers (placeholder implementations)
func healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
		"service": "ow-exporter",
		"version": "development",
	})
}

func indexHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"service": "ow-exporter",
		"version": "development",
		"status": "in development",
		"endpoints": map[string]string{
			"health": "/health",
			"metrics": "/metrics",
			"users": "/api/users",
		},
		"documentation": "https://github.com/lexfrei/tools/issues/439",
	})
}

func listUsersHandler(c echo.Context) error {
	// TODO: Implement user listing from SQLite
	return c.JSON(http.StatusOK, map[string]interface{}{
		"users": []string{},
		"total": 0,
		"status": "not implemented",
	})
}

func createUserHandler(c echo.Context) error {
	// TODO: Implement user creation
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"error": "not implemented",
		"message": "User creation will be implemented in next phase",
	})
}

func getUserHandler(c echo.Context) error {
	username := c.Param("username")
	// TODO: Implement user retrieval
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"error": "not implemented",
		"username": username,
		"message": "User retrieval will be implemented in next phase",
	})
}

func updateUserHandler(c echo.Context) error {
	username := c.Param("username")
	// TODO: Implement user update
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"error": "not implemented",
		"username": username,
		"message": "User update will be implemented in next phase",
	})
}

func deleteUserHandler(c echo.Context) error {
	username := c.Param("username")
	// TODO: Implement user deletion
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"error": "not implemented",
		"username": username,
		"message": "User deletion will be implemented in next phase",
	})
}

func metricsHandler(c echo.Context) error {
	// TODO: Implement Prometheus metrics
	return c.String(http.StatusOK, `# HELP ow_exporter_info Information about ow-exporter
# TYPE ow_exporter_info gauge
ow_exporter_info{version="development",status="in_development"} 1
`)
}
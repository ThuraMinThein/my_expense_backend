package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ThuraMinThein/my_expense_backend/config"
	"github.com/ThuraMinThein/my_expense_backend/db"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/handlers"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/repositories"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/routes"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/services"
	"github.com/ThuraMinThein/my_expense_backend/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/secure"
)

func main() {
	setupLogger()

	config.LoadConfig()

	migrateDatabase := false
	if err := db.DatabaseInit(migrateDatabase); err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middlewares.RequestIDMiddleware())
	r.Use(middlewares.LoggingMiddleware())
	r.Use(middlewares.ErrorHandlingMiddleware())

	secureMiddleware := secure.New(secure.Options{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
	})
	r.Use(func(c *gin.Context) {
		err := secureMiddleware.Process(c.Writer, c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "security middleware blocked request"})
			return
		}
		c.Next()
	})

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	//initialize repository, service, and handler
	repo := repositories.NewRepository(db.DB)
	services := services.NewServices(&repo)

	h := handlers.InitHandlers(services)

	routes.RegisterRoutes(r, h)

	startServer(r)
}

func setupLogger() {
	if os.Getenv("GIN_MODE") == "release" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func startServer(r *gin.Engine) {
	port := config.Config.ServerPort
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logrus.WithField("port", port).Info("Server has started")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("Server startup error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Warn("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("Forced shutdown: %v", err)
	}

	logrus.Info("Server exited cleanly")
}

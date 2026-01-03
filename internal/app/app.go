package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"project-mini-e-commerce/internal/config"
	"project-mini-e-commerce/internal/middleware"
	"project-mini-e-commerce/internal/routes"
	"project-mini-e-commerce/internal/validation"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Module interface {
	Routes() routes.Route
}

type Application struct {
	config *config.Config
	router *gin.Engine
}

func NewApplication(cfg *config.Config) *Application {
	if err := validation.InitValidator(); err != nil {
		log.Fatalf("Failed to initialize validator: %v", err)
	}

	go middleware.CleanupClient()
	r := gin.New()

	app := &Application{
		config: cfg,
		router: r,
	}

	app.registerModules()

	return app
}

func (a *Application) registerModules() {
	modules := []Module{
		NewUserModel(),
	}

	var moduleRoutes []routes.Route
	for _, m := range modules {
		moduleRoutes = append(moduleRoutes, m.Routes())
	}

	routes.RegisterRoutes(a.router, moduleRoutes...)
}

func (a *Application) Run() error {
	srv := &http.Server{
		Addr:    a.config.ServerAddress,
		Handler: a.router,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer stop()

	go func() {
		log.Printf("server started on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	log.Println("shutting down server...")

	srv.SetKeepAlivesEnabled(false)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	return nil
}

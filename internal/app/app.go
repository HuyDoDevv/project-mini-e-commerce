package app

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"project-mini-e-commerce/internal/config"
	"project-mini-e-commerce/internal/db"
	"project-mini-e-commerce/internal/db/sqlc"
	"project-mini-e-commerce/internal/routes"
	"project-mini-e-commerce/internal/validation"
	"project-mini-e-commerce/pkg/auth"
	"project-mini-e-commerce/pkg/cache"
	"project-mini-e-commerce/pkg/logger"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Module interface {
	Routes() routes.Route
}

type Application struct {
	config *config.Config
	router *gin.Engine
}

type ModuleContext struct {
	DB    sqlc.Querier
	Redis *redis.Client
}

func NewApplication(cfg *config.Config) *Application {
	if err := validation.InitValidator(); err != nil {
		logger.Logger.Fatal().Err(err).Msg("Failed to initialize validator")
	}
	r := gin.New()

	if err := db.InitDB(); err != nil {
		logger.Logger.Fatal().Err(err).Msg("Failed to initialize DB")
	}

	app := &Application{
		config: cfg,
		router: r,
	}

	app.registerModules()

	return app
}

func (a *Application) registerModules() {
	redisClient := config.NewRedisClient()
	cacheRedisService := cache.NewRedisCacheService(redisClient)
	tokenService := auth.NewJWTService(cacheRedisService)

	ctx := &ModuleContext{
		DB:    db.DB,
		Redis: redisClient,
	}

	modules := []Module{
		NewUserModel(ctx),
		NewAuthModule(ctx, tokenService, cacheRedisService),
	}

	var moduleRoutes []routes.Route
	for _, m := range modules {
		moduleRoutes = append(moduleRoutes, m.Routes())
	}

	routes.RegisterRoutes(a.router, tokenService, cacheRedisService, moduleRoutes...)
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
		logger.Logger.Info().Msgf("Server started on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.Logger.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	<-ctx.Done()

	logger.Logger.Info().Msg("Shutting down server...")

	srv.SetKeepAlivesEnabled(false)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	return nil
}

package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"project-mini-e-commerce/internal/config"
	"project-mini-e-commerce/internal/db"
	"project-mini-e-commerce/internal/db/sqlc"
	"project-mini-e-commerce/internal/routes"
	"project-mini-e-commerce/internal/validation"
	"project-mini-e-commerce/pkg/auth"
	"project-mini-e-commerce/pkg/cache"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
		log.Fatalf("Failed to initialize validator: %v", err)
	}
	loadEnv()
	r := gin.New()

	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
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
		NewAuthModule(ctx, tokenService),
	}

	var moduleRoutes []routes.Route
	for _, m := range modules {
		moduleRoutes = append(moduleRoutes, m.Routes())
	}

	routes.RegisterRoutes(a.router, tokenService, moduleRoutes...)
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

func loadEnv() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("No .env file found")
	}
}

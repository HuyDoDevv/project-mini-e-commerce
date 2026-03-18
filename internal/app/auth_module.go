package app

import (
	v1handler "project-mini-e-commerce/internal/handler/v1"
	"project-mini-e-commerce/internal/repository"
	"project-mini-e-commerce/internal/routes"
	v1routes "project-mini-e-commerce/internal/routes/v1"
	v1service "project-mini-e-commerce/internal/service/v1"
	"project-mini-e-commerce/pkg/auth"
	"project-mini-e-commerce/pkg/cache"
	"project-mini-e-commerce/pkg/mail"
	"project-mini-e-commerce/pkg/rabbitmq"
)

type AuthModule struct {
	routes routes.Route
}

func NewAuthModule(moduleCtx *ModuleContext, tokenService auth.TokenService, cacheService cache.RedisCacheService, mailService mail.EmailProviderService, rabbitmqService rabbitmq.RabbitMQService) *AuthModule {
	userRepo := repository.NewQueryUserRepository(moduleCtx.DB)

	authSer := v1service.NewAuthService(userRepo, tokenService, cacheService, mailService, rabbitmqService)

	authHand := v1handler.NewAuthHandler(authSer)

	authRoutes := v1routes.NewAuthRoutes(authHand)

	return &AuthModule{routes: authRoutes}
}

func (m *AuthModule) Routes() routes.Route {
	return m.routes
}

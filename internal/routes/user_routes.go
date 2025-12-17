package routes

import (
	"project-mini-e-commerce/internal/handler"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	handler *handler.UserHandler
}

func NewUserRoutes(handler *handler.UserHandler) *UserRoutes {
	return &UserRoutes{
		handler: handler,
	}
}

func (ur *UserRoutes) Register(r *gin.RouterGroup) {
	users := r.Group("/users")

	{
		users.GET("", ur.handler.GetAllUser)
		users.POST("", ur.handler.CreateUser)
		users.GET("/:uuid", ur.handler.GetByUserUUID)
		users.PUT("/:uuid", ur.handler.UpdateUser)
		users.DELETE("/:uuid", ur.handler.DeleteUser)
	}
}

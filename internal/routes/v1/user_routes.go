package v1routes

import (
	v1handler "project-mini-e-commerce/internal/handler/v1"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	handler *v1handler.UserHandler
}

func NewUserRoutes(handler *v1handler.UserHandler) *UserRoutes {
	return &UserRoutes{
		handler: handler,
	}
}

func (ur *UserRoutes) Register(r *gin.RouterGroup) {
	users := r.Group("/users")

	{
		users.GET("", ur.handler.GetAllUser)
		users.GET("/panic-users", ur.handler.PanicUser)
		users.POST("", ur.handler.CreateUser)
		users.GET("/:uuid", ur.handler.GetByUserUUID)
		users.PUT("/:uuid", ur.handler.UpdateUser)
		users.DELETE("/:uuid", ur.handler.DeleteUser)
	}
}

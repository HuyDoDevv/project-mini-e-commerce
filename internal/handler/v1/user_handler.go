package v1handler

import (
	v1service "project-mini-e-commerce/internal/service/v1"
)

type UserHandler struct {
	service v1service.UserService
}

func NewUserHandler(service v1service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) GetAllUser() {

}

func (uh *UserHandler) CreateUser() {

}

func (uh *UserHandler) GetByUserUUID() {

}
func (uh *UserHandler) UpdateUser() {

}
func (uh *UserHandler) DeleteUser() {

}

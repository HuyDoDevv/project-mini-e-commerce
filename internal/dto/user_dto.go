package dto

import "project-mini-e-commerce/internal/models"

type UserDTO struct {
	Uuid   string `json:"uuid"`
	Name   string `json:"full_name"`
	Email  string `json:"email"`
	Age    int    `json:"age"`
	Status string `json:"status"`
	Level  string `json:"level"`
}

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email,email_advanced"`
	Age      int    `json:"age" binding:"required,gt=0"`
	Password string `json:"password" binding:"required,password_strong"`
	Status   int    `json:"status" binding:"required,oneof=1 2"`
	Level    int    `json:"level" binding:"required,oneof=1 2"`
}
type UpdateUserInput struct {
	Name     string `json:"name" binding:"omitempty"`
	Email    string `json:"email" binding:"omitempty,email,email_advanced"`
	Age      int    `json:"age" binding:"omitempty,gt=0"`
	Password string `json:"password" binding:"omitempty,password_strong"`
	Status   int    `json:"status" binding:"omitempty,oneof=1 2"`
	Level    int    `json:"level" binding:"omitempty,oneof=1 2"`
}

func (input *CreateUserInput) MapCreateUserToModel() models.User {
	return models.User{
		Name:     input.Name,
		Email:    input.Email,
		Age:      input.Age,
		Password: input.Password,
		Status:   input.Status,
		Level:    input.Level,
	}
}

func (input *UpdateUserInput) MapUpdateUserToModel() models.User {
	return models.User{
		Name:     input.Name,
		Email:    input.Email,
		Age:      input.Age,
		Password: input.Password,
		Status:   input.Status,
		Level:    input.Level,
	}
}

func MapUserDTO(user models.User) *UserDTO {
	return &UserDTO{
		Uuid:   user.Uuid,
		Name:   user.Name,
		Age:    user.Age,
		Email:  user.Email,
		Status: MapStatusUser(user.Status),
		Level:  MapLevelUser(user.Level),
	}
}

func MapUsersDTO(users []models.User) []UserDTO {
	dtos := make([]UserDTO, 0, len(users))

	for _, user := range users {
		dtos = append(dtos, *MapUserDTO(user))
	}
	return dtos
}

func MapStatusUser(status int) string {
	switch status {
	case 1:
		return "Show"
	case 2:
		return "Hide"
	default:
		return "None"
	}
}
func MapLevelUser(status int) string {
	switch status {
	case 1:
		return "Admin"
	case 2:
		return "Client"
	default:
		return "None"
	}
}

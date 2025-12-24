package dto

import "project-mini-e-commerce/internal/models"

type UserDTO struct {
	Uuid   string `json:"uuid"`
	Name   string `json:"full_name"`
	Age    int    `json:"age"`
	Status string `json:"status"`
	Level  string `json:"level"`
}

func MapUserDTO(user models.User) *UserDTO {
	return &UserDTO{
		Uuid:   user.Uuid,
		Name:   user.Name,
		Age:    user.Age,
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

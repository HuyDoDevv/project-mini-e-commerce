package v1dto

type UserDTO struct {
	Uuid   string `json:"uuid"`
	Name   string `json:"full_name"`
	Age    int    `json:"age"`
	Status string `json:"status"`
	Level  string `json:"level"`
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

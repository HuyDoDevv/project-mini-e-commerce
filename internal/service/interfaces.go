package service

type UserService interface {
	GetAllUser()
	CreateUser()
	GetByUserUUID()
	UpdateUser()
	DeleteUser()
}

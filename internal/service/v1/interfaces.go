package v1service

type UserService interface {
	GetAllUser()
	CreateUser()
	GetByUserUUID()
	UpdateUser()
	DeleteUser()
}

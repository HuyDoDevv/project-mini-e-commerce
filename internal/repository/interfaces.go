package repository

type UserRepository interface {
	GetAll()
	Create()
	GetByUUID()
	Update()
	Delete()
}

package repository

type QueryUserRepository struct {
}

func NewQueryUserRepository() UserRepository {
	return &QueryUserRepository{}
}

func (ur *QueryUserRepository) GetAll()          {}
func (ur *QueryUserRepository) Create()          {}
func (ur *QueryUserRepository) GetByUUID()       {}
func (ur *QueryUserRepository) Update()          {}
func (ur *QueryUserRepository) Delete()          {}
func (ur *QueryUserRepository) FindUserByEmail() {}

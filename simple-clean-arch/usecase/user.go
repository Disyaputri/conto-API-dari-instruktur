package usecase

import (
	"simple-http/model"
	"simple-http/repository"
)

type UserUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) GetUsers() []model.User {
	return u.repo.FindAll()
}

func (u *UserUsecase) GetUserByID(id int) (model.User, error) {
	return u.repo.FindByID(id)
}

func (u *UserUsecase) CreateUser(user model.User) (model.User, error) {
	return u.repo.Create(user)
}

func (u *UserUsecase) UpdateUser(id int, request model.UserUpdateRequest) (model.User, error) {
	return u.repo.Update(id, request.Name)
}

func (u *UserUsecase) DeleteUser(id int) (model.User, error) {
	return u.repo.Delete(id)
}

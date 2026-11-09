package repository

import (
	"errors"
	"sync"

	"simple-http/model"
)

var (
	ErrNotFound    = errors.New("user not found")
	ErrDuplicateID = errors.New("user id already exists")
)

type UserRepository interface {
	FindAll() []model.User
	FindByID(id int) (model.User, error)
	Create(user model.User) (model.User, error)
	Update(id int, name string) (model.User, error)
	Delete(id int) (model.User, error)
}

type MemoryUserRepository struct {
	mu    sync.RWMutex
	users []model.User
}

func NewUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: []model.User{
			{ID: 1, Name: "John"},
			{ID: 2, Name: "Doe"},
		},
	}
}

func (r *MemoryUserRepository) FindAll() []model.User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]model.User, len(r.users))
	copy(users, r.users)
	return users
}

func (r *MemoryUserRepository) FindByID(id int) (model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}

	return model.User{}, ErrNotFound
}

func (r *MemoryUserRepository) Create(user model.User) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.users {
		if existing.ID == user.ID {
			return model.User{}, ErrDuplicateID
		}
	}

	r.users = append(r.users, user)
	return user, nil
}

func (r *MemoryUserRepository) Update(id int, name string) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.users {
		if r.users[i].ID == id {
			r.users[i].Name = name
			return r.users[i], nil
		}
	}

	return model.User{}, ErrNotFound
}

func (r *MemoryUserRepository) Delete(id int) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, user := range r.users {
		if user.ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return user, nil
		}
	}

	return model.User{}, ErrNotFound
}

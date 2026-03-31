package storage

import (
	"errors"
	"sync"

	"stocktrack-backend/internal/domain"
)

type InMemoryUserRepository struct {
	users map[string]*domain.User
	mu    sync.RWMutex
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *InMemoryUserRepository) Save(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user == nil || user.ID == "" {
		return errors.New("invalid user")
	}

	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) FindByEmail(email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (r *InMemoryUserRepository) FindByID(id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (r *InMemoryUserRepository) Exists(email string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if email == "" {
		return false
	}

	for _, user := range r.users {
		if user.Email == email {
			return true
		}
	}

	return false
}

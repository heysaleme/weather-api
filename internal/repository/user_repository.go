package repository

import (
	"context"
	"strings"
	"sync"

	"weather-api/internal/errs"
	"weather-api/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	List(ctx context.Context) ([]*model.User, error)
	Delete(ctx context.Context, id int64) error
}

type InMemoryUserRepository struct {
	mu      sync.RWMutex
	nextID  int64
	byID    map[int64]*model.User
	byEmail map[string]int64
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		nextID:  1,
		byID:    make(map[int64]*model.User),
		byEmail: make(map[string]int64),
	}
}

func (r *InMemoryUserRepository) Create(_ context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	emailKey := normalizeEmail(user.Email)
	if _, exists := r.byEmail[emailKey]; exists {
		return errs.Conflict("email already exists")
	}

	id := r.nextID
	r.nextID++

	userCopy := *user
	userCopy.ID = id

	r.byID[id] = &userCopy
	r.byEmail[emailKey] = id
	user.ID = id

	return nil
}

func (r *InMemoryUserRepository) GetByEmail(_ context.Context, email string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.byEmail[normalizeEmail(email)]
	if !ok {
		return nil, nil
	}

	user := *r.byID[id]
	return &user, nil
}

func (r *InMemoryUserRepository) GetByID(_ context.Context, id int64) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.byID[id]
	if !ok {
		return nil, nil
	}

	userCopy := *user
	return &userCopy, nil
}

func (r *InMemoryUserRepository) List(_ context.Context) ([]*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*model.User, 0, len(r.byID))
	for _, user := range r.byID {
		userCopy := *user
		result = append(result, &userCopy)
	}

	return result, nil
}

func (r *InMemoryUserRepository) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.byID[id]
	if !ok {
		return nil
	}

	delete(r.byEmail, normalizeEmail(user.Email))
	delete(r.byID, id)
	return nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

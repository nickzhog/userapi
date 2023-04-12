package filestorage

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nickzhog/userapi/internal/server/service/user"
	"github.com/nickzhog/userapi/pkg/logging"
)

var _ user.Repository = (*repository)(nil)

type repository struct {
	logger *logging.Logger

	file  *os.File
	list  map[string]user.User
	mutex *sync.RWMutex
}

func NewFileStorage(storepath string, logger *logging.Logger) *repository {
	file, err := os.OpenFile(storepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.Fatal(err)
	}

	list := make(map[string]user.User)
	json.NewDecoder(file).Decode(&list)

	return &repository{
		logger: logger,
		file:   file,
		list:   list,
		mutex:  new(sync.RWMutex),
	}
}

func (r *repository) Create(ctx context.Context, usr *user.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	usr.ID = uuid.New().String()
	usr.CreatedAt = time.Now()

	r.list[usr.ID] = *usr

	return r.saveUsersToFile()
}

func (r *repository) Update(ctx context.Context, id string, usr user.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	old, exist := r.list[id]
	if !exist {
		return user.ErrNotFound
	}

	usr.ID = id
	usr.CreatedAt = old.CreatedAt

	r.list[id] = usr

	return r.saveUsersToFile()
}

func (r *repository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exist := r.list[id]
	if !exist {
		return user.ErrNotFound
	}
	delete(r.list, id)

	return r.saveUsersToFile()
}

func (r *repository) FindOne(ctx context.Context, id string) (user.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	usr, exist := r.list[id]
	if !exist {
		return user.User{}, user.ErrNotFound
	}

	return usr, nil
}

func (r *repository) FindAll(ctx context.Context) ([]user.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	users := make([]user.User, 0, len(r.list))

	for _, u := range r.list {
		users = append(users, u)
	}

	return users, nil
}

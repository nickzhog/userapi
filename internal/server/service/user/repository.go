package user

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, id string, user *User) error
	Delete(ctx context.Context, id string) error
	FindOne(ctx context.Context, id string) (User, error)
	FindAll(ctx context.Context) ([]User, error)
}

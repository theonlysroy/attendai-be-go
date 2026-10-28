package storage

import (
	"context"

	"github.com/theonlysroy/attendai-be-go/internal/model"
)

type StudentStore interface {
	Create(ctx context.Context, s model.User) (model.User, error)
	GetByID(ctx context.Context, id int64) (model.User, error)
	List(ctx context.Context, limit int, offset int) ([]model.User, int64, error)
	UpdateByID(ctx context.Context, id int64) (model.User, error)
	DeleteByID(ctx context.Context, id int64) error
}

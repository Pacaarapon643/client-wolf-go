package port

import (
	"context"
)

type UserRepo interface {
	CountUser(ctx context.Context, count *int64) error
}

type UserService interface {
	CountUser(ctx context.Context, count *int64) error
}

package services

import (
	"context"
)

func (s *Service) CountUser(ctx context.Context, count *int64) error {
	return s.r.CountUser(ctx, count)
}

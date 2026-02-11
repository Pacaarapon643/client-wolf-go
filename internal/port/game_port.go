package port

import "context"

type GameService interface {
	RoleDistribution(ctx context.Context, roomId string) error
}

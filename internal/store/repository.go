package store

import (
	"context"

	"awesomeProject/internal/client"
	"awesomeProject/internal/reward"
)

// Repository defines all persistence operations.
// Implementations can be swapped (in-memory, SQL, etc).
type Repository interface {
	// Client operations
	CreateClient(ctx context.Context, name, email string) (*client.Client, error)
	ListClients(ctx context.Context) ([]client.Client, error)
	GetClient(ctx context.Context, id string) (*client.Client, error)
	AddAward(ctx context.Context, clientID string, awardType client.AwardType, pts int64) (*client.Award, error)
	GetAwards(ctx context.Context, clientID string) ([]client.Award, error)

	// Reward operations
	GetReward(ctx context.Context, id string) (*reward.Reward, error)
	ListRewards(ctx context.Context) ([]reward.Reward, error)
	Redeem(ctx context.Context, clientID, rewardID string, pointsCost int64) (int64, error)
}

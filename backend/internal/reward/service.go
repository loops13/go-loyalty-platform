package reward

import (
	"context"

	"GoLoyaltyPlatform/internal/client"
	"GoLoyaltyPlatform/internal/logging"
)

// Repository defines reward persistence operations.
type Repository interface {
	GetReward(ctx context.Context, id string) (*Reward, error)
	ListRewards(ctx context.Context) ([]Reward, error)
	Redeem(ctx context.Context, clientID, rewardID string, pointsCost int64) (int64, error)
}

// Service implements business logic for rewards.
type Service struct {
	repo          Repository
	clientService *client.Service
}

func NewService(repo Repository, clientService *client.Service) *Service {
	return &Service{
		repo:          repo,
		clientService: clientService,
	}
}

// Get fetches a reward by ID.
func (s *Service) Get(ctx context.Context, id string) (*Reward, error) {
	r, err := s.repo.GetReward(ctx, id)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, ErrNotFound
	}
	return r, nil
}

// List returns all available rewards.
func (s *Service) List(ctx context.Context) ([]Reward, error) {
	return s.repo.ListRewards(ctx)
}

// Redeem redeems a reward for a client.
// Business logic: verifies client exists, checks balance, deducts points.
func (s *Service) Redeem(ctx context.Context, clientID, rewardID string) (*Reward, int64, error) {
	logger := logging.FromContext(ctx)

	// Verify client exists
	c, err := s.clientService.Get(ctx, clientID)
	if err != nil {
		logger.Error("failed to load client for redemption", "client_id", clientID, "reward_id", rewardID, "error", err)
		return nil, 0, err
	}

	// Get reward
	reward, err := s.Get(ctx, rewardID)
	if err != nil {
		logger.Error("failed to load reward for redemption", "client_id", clientID, "reward_id", rewardID, "error", err)
		return nil, 0, err
	}

	// Check balance
	if c.PointBalance < reward.PointsCost {
		logger.Warn("insufficient points for redemption", "client_id", clientID, "reward_id", rewardID, "balance", c.PointBalance, "cost", reward.PointsCost)
		return nil, 0, client.ErrInsufficientPts
	}

	// Execute redemption
	newBalance, err := s.repo.Redeem(ctx, clientID, rewardID, reward.PointsCost)
	if err != nil {
		logger.Error("failed to redeem reward", "client_id", clientID, "reward_id", rewardID, "error", err)
		return nil, 0, err
	}

	logger.Info("reward redeemed", "client_id", clientID, "reward_id", rewardID, "balance", newBalance, "cost", reward.PointsCost)
	return reward, newBalance, nil
}

package client

import (
	"context"
	"errors"

	"awesomeProject/internal/logging"
)

// Repository defines client persistence operations.
type Repository interface {
	CreateClient(ctx context.Context, name, email string) (*Client, error)
	GetClient(ctx context.Context, id string) (*Client, error)
	AddAward(ctx context.Context, clientID string, awardType AwardType, pts int64) (*Award, error)
	GetAwards(ctx context.Context, clientID string) ([]Award, error)
}

// Service implements business logic for clients.
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new client.
func (s *Service) Create(ctx context.Context, name, email string) (*Client, error) {
	logger := logging.FromContext(ctx)

	c, err := s.repo.CreateClient(ctx, name, email)
	if err != nil {
		logger.Error("failed to create client", "error", err)
		return nil, err
	}

	logger.Info("client created", "client_id", c.ID)
	return c, nil
}

// Get fetches a client by ID.
func (s *Service) Get(ctx context.Context, id string) (*Client, error) {
	c, err := s.repo.GetClient(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrNotFound
	}
	return c, nil
}

// Award awards points for an action.
func (s *Service) Award(ctx context.Context, clientID, awardType string) (*Award, error) {
	logger := logging.FromContext(ctx)

	if !ValidAwardType(awardType) {
		logger.Warn("invalid award type", "client_id", clientID, "award_type", awardType)
		return nil, ErrInvalidAwardType
	}

	t := AwardType(awardType)
	pts := PointsForAward(t)

	award, err := s.repo.AddAward(ctx, clientID, t, pts)
	if err != nil {
		var ce *ClientError
		if errors.As(err, &ce) {
			logger.Warn("failed to award points", "client_id", clientID, "award_type", t, "error", err)
			return nil, err
		}
		logger.Error("failed to award points", "client_id", clientID, "award_type", t, "error", err)
		return nil, err
	}
	logger.Info("points awarded", "client_id", clientID, "award_type", t, "points", pts, "award_id", award.ID)
	return award, nil
}

// GetAwards retrieves award history for a client.
func (s *Service) GetAwards(ctx context.Context, clientID string) ([]Award, error) {
	// Verify client exists
	if _, err := s.Get(ctx, clientID); err != nil {
		return nil, err
	}
	return s.repo.GetAwards(ctx, clientID)
}

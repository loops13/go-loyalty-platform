package store

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"GoLoyaltyPlatform/internal/client"
	"GoLoyaltyPlatform/internal/reward"
)

// InMemoryStore implements Repository using in-memory maps.
type InMemoryStore struct {
	mu      sync.RWMutex
	clients map[string]*client.Client
	awards  map[string][]client.Award
	rewards map[string]*reward.Reward
}

// NewInMemoryStore creates and initializes an in-memory store.
func NewInMemoryStore() *InMemoryStore {
	store := &InMemoryStore{
		clients: make(map[string]*client.Client),
		awards:  make(map[string][]client.Award),
		rewards: make(map[string]*reward.Reward),
	}

	// Seed rewards
	store.rewards["r001"] = &reward.Reward{ID: "r001", Name: "Coffee Voucher", PointsCost: 500}
	store.rewards["r002"] = &reward.Reward{ID: "r002", Name: "Movie Ticket", PointsCost: 800}

	return store
}

// Client operations

func (s *InMemoryStore) CreateClient(ctx context.Context, name, email string) (*client.Client, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := genID("c")
	c := &client.Client{
		ID:           id,
		Name:         name,
		Email:        email,
		PointBalance: 0,
	}
	s.clients[id] = c
	return c, nil
}

func (s *InMemoryStore) GetClient(ctx context.Context, id string) (*client.Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.clients[id], nil
}

func (s *InMemoryStore) DeleteClient(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.clients[id]; !ok {
		return client.ErrNotFound
	}

	delete(s.clients, id)
	delete(s.awards, id)
	return nil
}

func (s *InMemoryStore) ListClients(ctx context.Context) ([]client.Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clients := make([]client.Client, 0, len(s.clients))
	for _, c := range s.clients {
		clients = append(clients, *c)
	}
	return clients, nil
}

func (s *InMemoryStore) AddAward(ctx context.Context, clientID string, awardType client.AwardType, pts int64) (*client.Award, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.clients[clientID]
	if !ok {
		return nil, client.ErrNotFound
	}

	award := &client.Award{
		ID:            genID("a"),
		ClientID:      clientID,
		Type:          awardType,
		PointsAwarded: pts,
		CreatedAt:     time.Now().UTC(),
	}

	c.PointBalance += pts
	s.awards[clientID] = append(s.awards[clientID], *award)

	return award, nil
}

func (s *InMemoryStore) GetAwards(ctx context.Context, clientID string) ([]client.Award, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.awards[clientID], nil
}

// Reward operations

func (s *InMemoryStore) GetReward(ctx context.Context, id string) (*reward.Reward, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.rewards[id], nil
}

func (s *InMemoryStore) ListRewards(ctx context.Context) ([]reward.Reward, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rewards := make([]reward.Reward, 0, len(s.rewards))
	for _, r := range s.rewards {
		rewards = append(rewards, *r)
	}
	return rewards, nil
}

func (s *InMemoryStore) Redeem(ctx context.Context, clientID, rewardID string, pointsCost int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.clients[clientID]
	if !ok {
		return 0, client.ErrNotFound
	}

	if c.PointBalance < pointsCost {
		return 0, client.ErrInsufficientPts
	}

	c.PointBalance -= pointsCost

	// Record redemption as an award with negative points
	award := client.Award{
		ID:            genID("a"),
		ClientID:      clientID,
		Type:          client.AwardType(fmt.Sprintf("REDEEM_%s", rewardID)),
		PointsAwarded: -pointsCost,
		CreatedAt:     time.Now().UTC(),
	}
	s.awards[clientID] = append(s.awards[clientID], award)

	return c.PointBalance, nil
}

func genID(prefix string) string {
	return fmt.Sprintf("%s_%d_%d", prefix, time.Now().UnixNano(), rand.Intn(1000))
}

package main

import (
	"sync"

	"github.com/Elbujito/2112/src/graphql-api/go/graph"
	"github.com/Elbujito/2112/src/graphql-api/go/graph/model"
	"github.com/go-redis/redis/v8"
)

// CustomResolver implements the root resolver for the GraphQL server.
type CustomResolver struct {
	rdb                   *redis.Client
	PositionSubscribers   map[string]map[string]chan *model.SatellitePosition
	VisibilitySubscribers map[string]chan []*model.SatelliteVisibility
	mu                    sync.Mutex
}

// NewCustomResolver initializes and returns a new CustomResolver.
func NewCustomResolver(redisClient *redis.Client) *CustomResolver {
	return &CustomResolver{
		rdb:                   redisClient,
		PositionSubscribers:   make(map[string]map[string]chan *model.SatellitePosition),
		VisibilitySubscribers: make(map[string]chan []*model.SatelliteVisibility),
	}
}

// Mutation returns the mutation resolver.
func (r *CustomResolver) Mutation() graph.MutationResolver {
	return &mutationResolver{CustomResolver: r}
}

// Query returns the query resolver.
func (r *CustomResolver) Query() graph.QueryResolver {
	return &queryResolver{CustomResolver: r}
}

// Subscription returns the subscription resolver.
func (r *CustomResolver) Subscription() graph.SubscriptionResolver {
	return &subscriptionResolver{CustomResolver: r}
}

// NotifyPositionSubscribers sends updates to position subscribers.
func (r *CustomResolver) NotifyPositionSubscribers(uid string, position *model.SatellitePosition) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if userChannels, ok := r.PositionSubscribers[uid]; ok {
		if ch, ok := userChannels[position.ID]; ok {
			ch <- position
		}
	}
}

// NotifyVisibilitySubscribers sends updates to visibility subscribers.
func (r *CustomResolver) NotifyVisibilitySubscribers(uid string, visibilities []*model.SatelliteVisibility) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ch, ok := r.VisibilitySubscribers[uid]; ok {
		ch <- visibilities
	}
}

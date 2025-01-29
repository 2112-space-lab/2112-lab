package main

import (
	"context"
	"fmt"

	"github.com/Elbujito/2112/src/graphql-api/go/graph/model"
)

type subscriptionResolver struct {
	*CustomResolver
}

// SatellitePositionUpdated resolves the subscription for real-time satellite position updates
func (s *subscriptionResolver) SatellitePositionUpdated(ctx context.Context, uid string, id string) (<-chan *model.SatellitePosition, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Ensure the user's position subscription map exists
	if _, exists := s.PositionSubscribers[uid]; !exists {
		s.PositionSubscribers[uid] = make(map[string]chan *model.SatellitePosition)
	}

	// Check if a subscription for this satellite already exists to prevent overwriting
	if _, exists := s.PositionSubscribers[uid][id]; exists {
		return nil, fmt.Errorf("subscription for satellite ID %s already exists for user %s", id, uid)
	}

	// Create a new channel for satellite position updates
	ch := make(chan *model.SatellitePosition, 1)
	s.PositionSubscribers[uid][id] = ch

	// Clean up subscription when the context is canceled
	go func() {
		<-ctx.Done()
		s.cleanupPositionSubscriber(uid, id, ch)
	}()

	return ch, nil
}

// SatelliteVisibilityUpdated resolves the subscription for real-time satellite visibility updates
func (s *subscriptionResolver) SatelliteVisibilityUpdated(ctx context.Context, uid string, userLocation model.UserLocationInput, startTime string, endTime string) (<-chan []*model.SatelliteVisibility, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if a subscription for visibility already exists to prevent overwriting
	if _, exists := s.VisibilitySubscribers[uid]; exists {
		// fmt.Errorf("visibility subscription already exists for user %s", uid)
		return nil, nil
	}

	// Create a new channel for visibility updates
	ch := make(chan []*model.SatelliteVisibility, 1)
	s.VisibilitySubscribers[uid] = ch

	// Clean up subscription when the context is canceled
	go func() {
		<-ctx.Done()
		s.cleanupVisibilitySubscriber(uid, ch)
	}()

	return ch, nil
}

// cleanupPositionSubscriber removes a position subscription and closes the channel
func (s *subscriptionResolver) cleanupPositionSubscriber(uid, id string, ch chan *model.SatellitePosition) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove the specific satellite subscription
	if userSubs, ok := s.PositionSubscribers[uid]; ok {
		if _, exists := userSubs[id]; exists {
			delete(userSubs, id)
			close(ch) // Close the channel to signal the end of the subscription
		}
		// If the user has no more subscriptions, clean up their map
		if len(userSubs) == 0 {
			delete(s.PositionSubscribers, uid)
		}
	}
}

// cleanupVisibilitySubscriber removes a visibility subscription and closes the channel
func (s *subscriptionResolver) cleanupVisibilitySubscriber(uid string, ch chan []*model.SatelliteVisibility) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove the visibility subscription for the user
	if _, exists := s.VisibilitySubscribers[uid]; exists {
		delete(s.VisibilitySubscribers, uid)
		close(ch) // Close the channel to signal the end of the subscription
	}
}

// Helper function to safely initialize a nested map if it doesn't exist
func (s *subscriptionResolver) ensurePositionSubscribers(uid string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.PositionSubscribers[uid]; !exists {
		s.PositionSubscribers[uid] = make(map[string]chan *model.SatellitePosition)
	}
}

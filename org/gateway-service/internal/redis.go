package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Elbujito/2112/src/graphql-api/go/graph/model"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

// Initialize Redis client
func initRedis(host, port string) {
	rdb = redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis successfully")
}

// Subscribe to Redis for position updates
func subscribeToRedisForPositionUpdates(resolver *CustomResolver) {
	pubsub := rdb.Subscribe(ctx, "satellite_positions")
	defer func() {
		if err := pubsub.Close(); err != nil {
			log.Printf("Error closing Redis subscription: %v", err)
		}
	}()

	log.Println("Subscribed to satellite_positions channel")

	for msg := range pubsub.Channel() {
		var position model.SatellitePosition
		if err := json.Unmarshal([]byte(msg.Payload), &position); err != nil {
			log.Printf("Error unmarshalling SatellitePosition: %v", err)
			continue
		}

		log.Printf("Received SatellitePosition update: %+v", position)
		// Notify relevant subscribers
		resolver.NotifyPositionSubscribers(position.ID, &position)
	}
}

// Subscribe to Redis for visibility updates
func subscribeToRedisForVisibilityUpdates(resolver *CustomResolver) {
	pubsub := rdb.Subscribe(ctx, "satellite_visibilities")
	defer func() {
		if err := pubsub.Close(); err != nil {
			log.Printf("Error closing Redis subscription: %v", err)
		}
	}()

	log.Println("Subscribed to satellite_visibilities channel")

	for msg := range pubsub.Channel() {
		var visibilities []*model.SatelliteVisibility
		if err := json.Unmarshal([]byte(msg.Payload), &visibilities); err != nil {
			log.Printf("Error unmarshalling SatelliteVisibility: %v", err)
			continue
		}

		log.Printf("Received SatelliteVisibility update for UID(s): %+v", extractUIDs(visibilities))
		// Notify relevant subscribers
		for _, visibility := range visibilities {
			resolver.NotifyVisibilitySubscribers(visibility.UserLocation.UID, visibilities)
		}
	}
}

// Helper function to extract UIDs from a list of SatelliteVisibility objects
func extractUIDs(visibilities []*model.SatelliteVisibility) []string {
	uids := make(map[string]struct{})
	for _, visibility := range visibilities {
		uids[visibility.UserLocation.UID] = struct{}{}
	}

	uidList := make([]string, 0, len(uids))
	for uid := range uids {
		uidList = append(uidList, uid)
	}
	return uidList
}

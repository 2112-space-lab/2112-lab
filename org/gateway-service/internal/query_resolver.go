package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Elbujito/2112/src/graphql-api/go/graph/model"
	"github.com/go-redis/redis/v8"
)

// queryResolver implements the GraphQL query resolver interface.
type queryResolver struct {
	*CustomResolver
}

// SatellitePosition retrieves the current position of a satellite by its unique ID.
func (q *queryResolver) SatellitePosition(ctx context.Context, id string) (*model.SatellitePosition, error) {
	// Fetch satellite position data from Redis
	data, err := q.CustomResolver.rdb.Get(ctx, "satellite_positions:"+id).Result()
	if err != nil {
		log.Printf("Error fetching SatellitePosition from Redis for ID %s: %v", id, err)
		return nil, fmt.Errorf("satellite position not found for ID %s: %w", id, err)
	}

	// Unmarshal JSON data into SatellitePosition struct
	var position model.SatellitePosition
	if err := json.Unmarshal([]byte(data), &position); err != nil {
		log.Printf("Error unmarshalling SatellitePosition for ID %s: %v", id, err)
		return nil, fmt.Errorf("failed to parse satellite position data for ID %s: %w", id, err)
	}

	return &position, nil
}

// SatelliteTle retrieves the TLE (Two-Line Element) data of a satellite by its unique ID.
func (q *queryResolver) SatelliteTle(ctx context.Context, id string) (*model.SatelliteTle, error) {
	// Fetch TLE data from Redis
	data, err := q.CustomResolver.rdb.Get(ctx, "satellite_tle:"+id).Result()
	if err != nil {
		log.Printf("Error fetching SatelliteTle from Redis for ID %s: %v", id, err)
		return nil, fmt.Errorf("TLE data not found for satellite ID %s: %w", id, err)
	}

	// Unmarshal JSON data into SatelliteTle struct
	var tle model.SatelliteTle
	if err := json.Unmarshal([]byte(data), &tle); err != nil {
		log.Printf("Error unmarshalling SatelliteTle for ID %s: %v", id, err)
		return nil, fmt.Errorf("failed to parse TLE data for ID %s: %w", id, err)
	}

	return &tle, nil
}

// SatellitePositionsInRange retrieves the positions of a satellite within a specific time range.
func (q *queryResolver) SatellitePositionsInRange(ctx context.Context, id string, startTime string, endTime string) ([]*model.SatellitePosition, error) {
	// Parse startTime and endTime into time.Time objects
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse startTime %s: %w", startTime, err)
	}
	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse endTime %s: %w", endTime, err)
	}

	// Prepare Redis query keys and timestamps
	key := fmt.Sprintf("satellite_positions:%s", id)
	startTimestamp := strconv.FormatInt(start.Unix(), 10)
	endTimestamp := strconv.FormatInt(end.Unix(), 10)

	// Query Redis for position data within the specified range
	results, err := q.CustomResolver.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: startTimestamp,
		Max: endTimestamp,
	}).Result()
	if err != nil {
		log.Printf("Error querying Redis for satellite positions in range for ID %s: %v", id, err)
		return nil, fmt.Errorf("failed to query Redis for satellite positions for ID %s: %w", id, err)
	}

	// Parse results into SatellitePosition objects
	var positions []*model.SatellitePosition
	for _, raw := range results {
		var position model.SatellitePosition
		if err := json.Unmarshal([]byte(raw), &position); err != nil {
			log.Printf("Error unmarshalling SatellitePosition for ID %s: %v", id, err)
			continue
		}
		positions = append(positions, &position)
	}

	return positions, nil
}

func (q *queryResolver) CachedSatelliteVisibilities(ctx context.Context, uid string, userLocation model.UserLocationInput, startTime string, endTime string) ([]*model.SatelliteVisibility, error) {
	// Prepare Redis key for cached visibilities
	key := fmt.Sprintf("satellite_visibilities:%s", uid)

	// Fetch cached visibility data from Redis
	data, err := q.CustomResolver.rdb.Get(ctx, key).Result()
	if err != nil {
		log.Printf("Error fetching cached satellite visibilities for UID %s: %v", uid, err)
		return nil, fmt.Errorf("cached visibilities not found for UID %s: %w", uid, err) 
	}

	// Parse cached visibility data into SatelliteVisibility objects
	var visibilities []*model.SatelliteVisibility
	if err := json.Unmarshal([]byte(data), &visibilities); err != nil {
		log.Printf("Error unmarshalling CachedSatelliteVisibilities for UID %s: %v", uid, err)
		log.Printf("Payload causing the error: %s", data) // Log the raw payload
		return nil, fmt.Errorf("failed to parse cached satellite visibilities for UID %s: %w", uid, err)
	}

	// Deduplicate the visibilities
	distinctVisibilities := distinctSatelliteVisibilities(visibilities)

	return distinctVisibilities, nil
}

// Helper function to remove duplicates
func distinctSatelliteVisibilities(visibilities []*model.SatelliteVisibility) []*model.SatelliteVisibility {
	seen := make(map[string]bool) // Use a map to track unique entries
	var distinct []*model.SatelliteVisibility

	for _, visibility := range visibilities {
		// Create a unique key for each SatelliteVisibility (e.g., combination of important fields)
		key := fmt.Sprintf("%s_%s_%s", visibility.SatelliteID, visibility.Aos, visibility.Los)
		if !seen[key] {
			seen[key] = true
			distinct = append(distinct, visibility)
		}
	}

	return distinct
}

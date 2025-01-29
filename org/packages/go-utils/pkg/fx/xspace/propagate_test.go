package xspace

import (
	"math"
	"testing"
	"time"

	xconstants "github.com/org/2112-space-lab/org/go-utils/pkg/fx/xconstants"
)

// Mock data for testing
const (
	mockTLELine1 = "1 25544U 98067A   21275.91835648  .00002907  00000-0  58234-4 0  9992"
	mockTLELine2 = "2 25544  51.6442 176.8457 0003392  45.8666  36.0921 15.48815362312356"
)

func TestPropagateRange(t *testing.T) {
	// Define the start and end times and the interval
	startTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, time.January, 1, 1, 0, 0, 0, time.UTC)
	interval := 10 * time.Minute

	// Call the function under test
	positions, err := PropagateRange(mockTLELine1, mockTLELine2, startTime, endTime, interval)

	// Check for unexpected errors
	if err != nil {
		t.Fatalf("PropagateRange returned an error: %v", err)
	}

	// Validate the number of positions generated
	expectedCount := int(endTime.Sub(startTime)/interval) + 1
	if len(positions) != expectedCount {
		t.Errorf("Expected %d positions, got %d", expectedCount, len(positions))
	}

	// Define expected values for the first and last positions (pre-calculated)
	expectedFirstPosition := SatellitePosition{
		Latitude:  29.463023839407164,
		Longitude: 6.630689954904137,
		Altitude:  413.82619943027055, // Altitude in kilometers
		Time:      startTime,
	}

	expectedLastPosition := SatellitePosition{
		Latitude:  -51.77103753887528,
		Longitude: -122.41576121987946,
		Altitude:  426.890472040785, // Altitude in kilometers
		Time:      endTime,
	}

	// Define a tolerance in kilometers
	tolerance := 36.0

	// Helper function to calculate distance in kilometers for latitude/longitude
	latLonToKm := func(lat1, lon1, lat2, lon2 float64) float64 {

		dLat := (lat2 - lat1) * xconstants.PI_DIVIDE_BY_180
		dLon := (lon2 - lon1) * xconstants.PI_DIVIDE_BY_180
		lat1 = lat1 * xconstants.PI_DIVIDE_BY_180
		lat2 = lat2 * xconstants.PI_DIVIDE_BY_180

		a := math.Sin(dLat/2)*math.Sin(dLat/2) +
			math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
		c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
		return xconstants.EARTH_RADIUS_KM * c
	}

	// Validate the first position
	firstPosition := positions[0]
	distanceFirst := latLonToKm(firstPosition.Latitude, firstPosition.Longitude, expectedFirstPosition.Latitude, expectedFirstPosition.Longitude)
	if distanceFirst > tolerance {
		t.Errorf("First position mismatch: distance exceeds tolerance. Expected (%f, %f), got (%f, %f), distance: %f km",
			expectedFirstPosition.Latitude, expectedFirstPosition.Longitude,
			firstPosition.Latitude, firstPosition.Longitude, distanceFirst)
	}
	if math.Abs(firstPosition.Altitude-expectedFirstPosition.Altitude) > tolerance {
		t.Errorf("First position altitude mismatch: expected %f, got %f", expectedFirstPosition.Altitude, firstPosition.Altitude)
	}

	// Validate the last position
	lastPosition := positions[len(positions)-1]
	distanceLast := latLonToKm(lastPosition.Latitude, lastPosition.Longitude, expectedLastPosition.Latitude, expectedLastPosition.Longitude)
	if distanceLast > tolerance {
		t.Errorf("Last position mismatch: distance exceeds tolerance. Expected (%f, %f), got (%f, %f), distance: %f km",
			expectedLastPosition.Latitude, expectedLastPosition.Longitude,
			lastPosition.Latitude, lastPosition.Longitude, distanceLast)
	}
	if math.Abs(lastPosition.Altitude-expectedLastPosition.Altitude) > tolerance {
		t.Errorf("Last position altitude mismatch: expected %f, got %f", expectedLastPosition.Altitude, lastPosition.Altitude)
	}
}

func TestPropagateValidateEachPosition(t *testing.T) {
	startTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, time.January, 1, 1, 0, 0, 0, time.UTC)
	interval := 10 * time.Minute

	positions, err := PropagateRange(mockTLELine1, mockTLELine2, startTime, endTime, interval)
	if err != nil {
		t.Fatalf("PropagateRange returned an error: %v", err)
	}

	expectedPositions := []SatellitePosition{
		{Latitude: 29.463023839407164, Longitude: 6.630689954904137, Altitude: 413.82619943027055},
		{Latitude: 50.11705197638757, Longitude: 48.037827753096316, Altitude: 418.90763849077905},
		{Latitude: 44.805945200576176, Longitude: 103.74490357389301, Altitude: 417.6752237256438},
		{Latitude: 19.305749788784198, Longitude: 136.61045112910892, Altitude: 412.56012381762275},
		{Latitude: -11.02789372404546, Longitude: 158.89245653528198, Altitude: 413.7518989253967},
		{Latitude: -38.8388318075352, Longitude: -173.15539299396244, Altitude: 422.20416801792356},
		{Latitude: -51.77103753887528, Longitude: -122.41576121987946, Altitude: 426.890472040785},
	}

	// Validate the number of positions generated
	if len(positions) != len(expectedPositions) {
		t.Errorf("Expected %d positions, got %d", len(expectedPositions), len(positions))
	}

	// Define a tolerance in kilometers
	toleranceKm := 40.0

	// Helper function to calculate distance in kilometers for latitude/longitude
	latLonToKm := func(lat1, lon1, lat2, lon2 float64) float64 {
		dLat := (lat2 - lat1) * xconstants.PI_DIVIDE_BY_180
		dLon := (lon2 - lon1) * xconstants.PI_DIVIDE_BY_180
		lat1 = lat1 * xconstants.PI_DIVIDE_BY_180
		lat2 = lat2 * xconstants.PI_DIVIDE_BY_180

		a := math.Sin(dLat/2)*math.Sin(dLat/2) +
			math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
		c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
		return xconstants.EARTH_RADIUS_KM * c
	}

	// Validate each position
	for i, position := range positions {
		if i >= len(expectedPositions) {
			break
		}

		expected := expectedPositions[i]
		distance := latLonToKm(position.Latitude, position.Longitude, expected.Latitude, expected.Longitude)
		if distance > toleranceKm {
			t.Errorf("Position %d mismatch: distance exceeds tolerance. Expected (%f, %f), got (%f, %f), distance: %f km",
				i, expected.Latitude, expected.Longitude, position.Latitude, position.Longitude, distance)
		}
		if math.Abs(position.Altitude-expected.Altitude) > toleranceKm {
			t.Errorf("Position %d altitude mismatch: expected %f, got %f", i, expected.Altitude, position.Altitude)
		}
	}
}

package xspace

import (
	"math"
	"testing"
	"time"

	"github.com/joshuaferrara/go-satellite"
	xpolygon "github.com/org/2112-space-lab/org/go-utils/pkg/fx/xpolygon"
)

func TestHaversineDistance(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lon1     float64
		lat2     float64
		lon2     float64
		expected float64
	}{
		{
			name:     "Same point",
			lat1:     0.0,
			lon1:     0.0,
			lat2:     0.0,
			lon2:     0.0,
			expected: 0.0,
		},
		{
			name:     "Points on the equator",
			lat1:     0.0,
			lon1:     0.0,
			lat2:     0.0,
			lon2:     90.0,
			expected: 10007.54, // Approximate distance
		},
		{
			name:     "Points on the same meridian",
			lat1:     0.0,
			lon1:     0.0,
			lat2:     90.0,
			lon2:     0.0,
			expected: 10007.54, // Approximate distance
		},
		{
			name:     "Points with small distance",
			lat1:     52.2296756,
			lon1:     21.0122287,
			lat2:     52.406374,
			lon2:     16.9251681,
			expected: 278.546, // Approximate distance
		},
		{
			name:     "Points with large distance",
			lat1:     36.12,
			lon1:     -86.67,
			lat2:     33.94,
			lon2:     -118.40,
			expected: 2886.07, // Approximate distance
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HaversineDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2, 0, 0)
			if math.Abs(result-tt.expected) > 1.0 { // Allow a small margin of error
				t.Errorf("HaversineDistance() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIntersectsEdge(t *testing.T) {
	tests := []struct {
		name         string
		tileCenter   xpolygon.Point // Replacing edge with center point
		satellitePos xpolygon.Point
		tileRadiusKm float64
		altitude     float64
		expected     bool
	}{
		{
			name:         "Satellite within radius of tile center (start point)",
			tileCenter:   xpolygon.Point{Latitude: 0.0, Longitude: 0.0}, // Center of the tile is the start point
			satellitePos: xpolygon.Point{Latitude: 0.0, Longitude: 0.0},
			tileRadiusKm: 1.0,
			altitude:     0.0,
			expected:     true,
		},
		{
			name:         "Satellite within radius of tile center (end point)",
			tileCenter:   xpolygon.Point{Latitude: 0.5, Longitude: 0.5}, // Center of the tile is the midpoint between start and end
			satellitePos: xpolygon.Point{Latitude: 0.5, Longitude: 0.5},
			tileRadiusKm: 1.0,
			altitude:     0.0,
			expected:     true,
		},
		{
			name:         "Satellite outside radius of tile center",
			tileCenter:   xpolygon.Point{Latitude: 0.5, Longitude: 0.5},
			satellitePos: xpolygon.Point{Latitude: 2.0, Longitude: 2.0},
			tileRadiusKm: 1.0,
			altitude:     0.0,
			expected:     false,
		},
		{
			name:         "Satellite exactly 1 km away in latitude direction",
			tileCenter:   xpolygon.Point{Latitude: 0.0, Longitude: 0.0},
			satellitePos: xpolygon.Point{Latitude: 0.00899322, Longitude: 0.0}, // Exactly 1 km away along latitude
			tileRadiusKm: 1.0,
			altitude:     0.0,
			expected:     true,
		},
		{
			name:         "Satellite exactly 1 km away in longitude direction",
			tileCenter:   xpolygon.Point{Latitude: 0.0, Longitude: 0.0},
			satellitePos: xpolygon.Point{Latitude: 0.0, Longitude: 0.00899322}, // Exactly 1 km away along longitude
			tileRadiusKm: 1.0,
			altitude:     0.0,
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Calculate distance from satellite to the center of the tile
			centerDistance := HaversineDistance(tt.satellitePos.Latitude, tt.satellitePos.Longitude, tt.tileCenter.Latitude, tt.tileCenter.Longitude, 0, 0)

			// Log detailed information about the distance calculation
			t.Logf("[DEBUG] Satellite Position: (%.6f, %.6f)", tt.satellitePos.Latitude, tt.satellitePos.Longitude)
			t.Logf("[DEBUG] Tile Center: (%.6f, %.6f)", tt.tileCenter.Latitude, tt.tileCenter.Longitude)
			t.Logf("[DEBUG] Calculated Center Distance: %.6f km", centerDistance)
			t.Logf("[DEBUG] Tile Radius: %.6f km", tt.tileRadiusKm)

			// Adding margin of error for comparison
			marginOfError := 0.01
			t.Logf("[DEBUG] Margin of Error: %.6f km", marginOfError)

			// Check if satellite is within the tile radius from the center of the tile
			if got := Intersects(tt.tileCenter, tt.satellitePos, tt.tileRadiusKm, tt.altitude); got != tt.expected {
				t.Errorf("IntersectsEdge() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestComputeAOS(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(1 * time.Hour) // Extended time window
	timeStep := 5 * time.Second             // Reduced time step
	tileRadiusKm := 7000.0                  // Increased tile radius to 7000 km

	tleLine1 := "1 25544U 98067A   20344.54791435  .00001234  00000-0  29746-4 0  9998"
	tleLine2 := "2 25544  51.6456 212.9669 0001235 341.2074 106.3520 15.48921140255678"
	satrec := satellite.TLEToSat(tleLine1, tleLine2, satellite.GravityWGS84)

	tests := []struct {
		name          string
		vertices      []xpolygon.Point
		expectedAOS   bool
		expectedMaxEl bool
	}{
		{
			name: "Low Latitude Tile",
			vertices: []xpolygon.Point{
				{Latitude: 0.0, Longitude: 0.0},
				{Latitude: 0.0, Longitude: 10.0},
				{Latitude: 10.0, Longitude: 0.0},
				{Latitude: 10.0, Longitude: 10.0},
			},
			expectedAOS:   true,
			expectedMaxEl: true,
		},
		{
			name: "High Latitude Tile",
			vertices: []xpolygon.Point{
				{Latitude: 45.0, Longitude: 45.0},
				{Latitude: 45.0, Longitude: 55.0},
				{Latitude: 55.0, Longitude: 45.0},
				{Latitude: 55.0, Longitude: 55.0},
			},
			expectedAOS:   true,
			expectedMaxEl: true,
		},
		{
			name: "Out of Range Tile",
			vertices: []xpolygon.Point{
				{Latitude: 80.0, Longitude: 80.0},
				{Latitude: 80.0, Longitude: 90.0},
				{Latitude: 90.0, Longitude: 80.0},
				{Latitude: 90.0, Longitude: 90.0},
			},
			expectedAOS:   false,
			expectedMaxEl: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var maxElevation float64 = -1.0

			aos := ComputeAOS(satrec, tt.vertices[0], tileRadiusKm, startTime, endTime, timeStep, &maxElevation)

			if tt.expectedAOS && aos.IsZero() {
				t.Errorf("Expected AOS but got none.")
			} else if !tt.expectedAOS && !aos.IsZero() {
				t.Errorf("Did not expect AOS but got one.")
			}

			if tt.expectedMaxEl && maxElevation <= 0 {
				t.Errorf("Expected positive MaxElevation but got %.2f", maxElevation)
			} else if !tt.expectedMaxEl && maxElevation > 0 {
				t.Errorf("Did not expect MaxElevation but got %.2f", maxElevation)
			}
		})
	}
}

func TestCalculateIntegratedElevation(t *testing.T) {
	// Define test cases with satellite position, altitude, ground point, and expected elevation range
	testCases := []struct {
		name           string
		satellitePos   xpolygon.Point
		satelliteAltKm float64
		groundPoint    xpolygon.Point
	}{
		{
			name:           "Directly Overhead",
			satellitePos:   xpolygon.Point{Latitude: 0.0, Longitude: 0.0},
			satelliteAltKm: 500.0,
			groundPoint:    xpolygon.Point{Latitude: 0.0, Longitude: 0.0},
		},
		{
			name:           "Near Horizon",
			satellitePos:   xpolygon.Point{Latitude: 10.0, Longitude: 10.0},
			satelliteAltKm: 500.0,
			groundPoint:    xpolygon.Point{Latitude: 10.0, Longitude: 20.0},
		},
		{
			name:           "Far Distance",
			satellitePos:   xpolygon.Point{Latitude: 45.0, Longitude: 45.0},
			satelliteAltKm: 1000.0,
			groundPoint:    xpolygon.Point{Latitude: -45.0, Longitude: -45.0},
		},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			elevation := CalculateIntegratedElevationFromPoint(tc.satellitePos, tc.satelliteAltKm, tc.groundPoint)
			if elevation <= 0 {
				t.Errorf("For %s, expected positive elevation but got %.2f", tc.name, elevation)
			} else {
				t.Logf("For %s, calculated elevation: %.2f degrees", tc.name, elevation)
			}
		})
	}
}

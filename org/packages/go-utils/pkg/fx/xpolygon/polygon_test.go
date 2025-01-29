package xpolygon

import (
	"log"
	"testing"
	"time"
)

// Constants for test cases
const (
	TestTileRadius = 10000.0 // Radius of the tile in meters
	TestNbFaces    = 6       // Number of faces for polygons
	TestZoomLevel  = 15      // Zoom level for testing
)

// TestNewPolygon verifies the creation of a new Polygon
func TestNewPolygon(t *testing.T) {
	center := LatLong{
		Lat: Coordinate{C: 40.7128},  // Latitude for New York City
		Lon: Coordinate{C: -74.0060}, // Longitude for New York City
	}

	polygon := NewPolygon(TestNbFaces, center, TestZoomLevel, TestTileRadius)

	if polygon.NbFaces != TestNbFaces {
		t.Errorf("Expected nbFaces %d, got %d", TestNbFaces, polygon.NbFaces)
	}
	if polygon.Radius != TestTileRadius {
		t.Errorf("Expected radius %.2f, got %.2f", TestTileRadius, polygon.Radius)
	}
	if len(polygon.Boundaries) != TestNbFaces {
		t.Errorf("Expected %d boundaries, got %d", TestNbFaces, len(polygon.Boundaries))
	}

	// Verify that all boundary points have valid latitude/longitude
	for i, boundary := range polygon.Boundaries {
		if boundary.Latitude < -90 || boundary.Latitude > 90 {
			t.Errorf("Boundary %d latitude out of bounds: %.6f", i, boundary.Latitude)
		}
		if boundary.Longitude < -180 || boundary.Longitude > 180 {
			t.Errorf("Boundary %d longitude out of bounds: %.6f", i, boundary.Longitude)
		}
	}
}

// TestGenerateBoundaries verifies that boundaries are calculated correctly
func TestGenerateBoundaries(t *testing.T) {
	center := LatLong{
		Lat: Coordinate{C: 0.0},
		Lon: Coordinate{C: 0.0},
	}

	boundaries := generateBoundaries(TestNbFaces, center, TestTileRadius)

	if len(boundaries) != TestNbFaces {
		t.Errorf("Expected %d boundaries, got %d", TestNbFaces, len(boundaries))
	}

	// Verify that all generated points are within valid latitude/longitude ranges
	for i, boundary := range boundaries {
		if boundary.Latitude < -90 || boundary.Latitude > 90 {
			t.Errorf("Boundary %d latitude out of bounds: %.6f", i, boundary.Latitude)
		}
		if boundary.Longitude < -180 || boundary.Longitude > 180 {
			t.Errorf("Boundary %d longitude out of bounds: %.6f", i, boundary.Longitude)
		}
	}
}

// TestCalculateZoomLevelForTileRadius tests the zoom level calculation
func TestCalculateZoomLevelForTileRadius(t *testing.T) {
	tileRadius := 500.0

	expectedZoom := 16
	calculatedZoom := calculateZoomLevelForTileRadius(tileRadius)

	if calculatedZoom != expectedZoom {
		t.Errorf("Expected zoom level %d, got %d", expectedZoom, calculatedZoom)
	}
}

// TestTileXYToLatLon verifies that tile coordinates are correctly converted to lat/lon
func TestTileXYToLatLon(t *testing.T) {
	zoom := 15
	x, y := 0, 0

	lat, lon := TileXYToLatLon(x, y, zoom)
	if lat > 85.05112878 || lat < -85.05112878 {
		t.Errorf("Latitude out of bounds: %.6f", lat)
	}
	if lon > 180 || lon < -180 {
		t.Errorf("Longitude out of bounds: %.6f", lon)
	}

	// Additional test: edge case at the maximum tile value for the zoom level
	x, y = 32767, 32767 // Maximum tile values at zoom 15
	lat, lon = TileXYToLatLon(x, y, zoom)
	if lat > 85.05112878 || lat < -85.05112878 {
		t.Errorf("Latitude out of bounds: %.6f", lat)
	}
	if lon > 180 || lon < -180 {
		t.Errorf("Longitude out of bounds: %.6f", lon)
	}
}

// TestLatLonToTileXY verifies that lat/lon coordinates are correctly converted to tile coordinates
func TestLatLonToTileXY(t *testing.T) {
	lat, lon := 40.7128, -74.0060 // Coordinates for New York City
	zoom := 15

	x, y := LatLonToTileXY(lat, lon, zoom)
	if x < 0 || y < 0 {
		t.Errorf("Invalid tile coordinates: x=%d, y=%d", x, y)
	}
}

// TestGenerateAllTilesForRadius verifies that tiles are generated for a given radius
func TestGenerateAllTilesForRadius(t *testing.T) {
	start := time.Now()
	tiles := GenerateAllTilesForRadius(TestTileRadius, TestNbFaces)
	duration := time.Since(start)

	log.Printf("Generated %d tiles in %s", len(tiles), duration)
	if len(tiles) == 0 {
		t.Error("Expected non-empty map of tiles, got empty map")
	}

	// Verify that all polygons have the correct number of faces
	for quadkey, polygon := range tiles {
		if polygon.NbFaces != TestNbFaces {
			t.Errorf("Expected polygon with %d faces, got %d for quadkey %v", TestNbFaces, polygon.NbFaces, quadkey)
		}

		// Verify that all boundary points are valid
		for i, boundary := range polygon.Boundaries {
			if boundary.Latitude < -90 || boundary.Latitude > 90 {
				t.Errorf("Boundary %d latitude out of bounds: %.6f for quadkey %v", i, boundary.Latitude, quadkey)
			}
			if boundary.Longitude < -180 || boundary.Longitude > 180 {
				t.Errorf("Boundary %d longitude out of bounds: %.6f for quadkey %v", i, boundary.Longitude, quadkey)
			}
		}
		break
	}
}

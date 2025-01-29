package xspace

import (
	"math"
	"testing"
)

// TestLatLonToCartesian tests the conversion of lat/lon/altitude to Cartesian coordinates
func TestLatLonToCartesian(t *testing.T) {
	tests := []struct {
		name      string
		latitude  float64
		longitude float64
		altitude  float64 // altitude in kilometers
		expectedX float64
		expectedY float64
		expectedZ float64
	}{
		{
			name:      "High altitude point",
			latitude:  45.0,
			longitude: 45.0,
			altitude:  1.0,     // altitude in kilometers
			expectedX: 3186.00, // Expected value of X after conversion
			expectedY: 3186.00, // Expected value of Y after conversion
			expectedZ: 4505.68, // Expected value of Z after conversion
		},
		{
			name:      "Low altitude point",
			latitude:  0.0,
			longitude: 0.0,
			altitude:  0.5,     // altitude in kilometers
			expectedX: 6371.50, // Expected value of X after conversion
			expectedY: 0.00,    // Expected value of Y after conversion
			expectedZ: 0.00,    // Expected value of Z after conversion
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function to test
			x, y, z := LatLonToCartesian(tt.latitude, tt.longitude, tt.altitude)

			// Check if the calculated values are within an acceptable margin of error
			if !almostEqual(x, tt.expectedX, 0.01) {
				t.Errorf("Expected X: %.2f, but got X: %.2f", tt.expectedX, x)
			}
			if !almostEqual(y, tt.expectedY, 0.01) {
				t.Errorf("Expected Y: %.2f, but got Y: %.2f", tt.expectedY, y)
			}
			if !almostEqual(z, tt.expectedZ, 0.01) {
				t.Errorf("Expected Z: %.2f, but got Z: %.2f", tt.expectedZ, z)
			}

			// Log the final result for debugging
			t.Logf("Calculated Cartesian Coordinates for %v: X=%.2f, Y=%.2f, Z=%.2f", tt.name, x, y, z)
		})
	}
}

// almostEqual compares two floats to a given precision (tolerance)
func almostEqual(a, b, tolerance float64) bool {
	return (a-b) < tolerance && (b-a) < tolerance
}

// Test for Normalize function
func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		x, y, z  float64
		expected [3]float64
	}{
		{
			name:     "Normalize unit vector",
			x:        1.0,
			y:        0.0,
			z:        0.0,
			expected: [3]float64{1.0, 0.0, 0.0},
		},
		{
			name:     "Normalize non-unit vector",
			x:        3.0,
			y:        4.0,
			z:        0.0,
			expected: [3]float64{0.6, 0.8, 0.0}, // (3, 4, 0) has magnitude 5, so divide by 5
		},
		{
			name:     "Normalize negative vector",
			x:        -1.0,
			y:        -2.0,
			z:        -2.0,
			expected: [3]float64{-0.333, -0.667, -0.667}, // Normalized vector for (-1, -2, -2)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y, z := Normalize(tt.x, tt.y, tt.z)
			if math.Abs(x-tt.expected[0]) > 0.01 || math.Abs(y-tt.expected[1]) > 0.01 || math.Abs(z-tt.expected[2]) > 0.01 {
				t.Errorf("Normalize() = (%.2f, %.2f, %.2f), want (%.2f, %.2f, %.2f)", x, y, z, tt.expected[0], tt.expected[1], tt.expected[2])
			}
		})
	}
}

// Test for DotProduct function
func TestDotProduct(t *testing.T) {
	tests := []struct {
		name                   string
		x1, y1, z1, x2, y2, z2 float64
		expected               float64
	}{
		{
			name: "Dot product of same direction vectors",
			x1:   1.0, y1: 1.0, z1: 1.0,
			x2: 1.0, y2: 1.0, z2: 1.0,
			expected: 3.0, // 1*1 + 1*1 + 1*1 = 3
		},
		{
			name: "Dot product of perpendicular vectors",
			x1:   1.0, y1: 0.0, z1: 0.0,
			x2: 0.0, y2: 1.0, z2: 0.0,
			expected: 0.0, // 1*0 + 0*1 + 0*0 = 0
		},
		{
			name: "Dot product of opposite vectors",
			x1:   1.0, y1: 0.0, z1: 0.0,
			x2: -1.0, y2: 0.0, z2: 0.0,
			expected: -1.0, // 1*(-1) + 0*0 + 0*0 = -1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DotProduct(tt.x1, tt.y1, tt.z1, tt.x2, tt.y2, tt.z2)
			if result != tt.expected {
				t.Errorf("DotProduct() = %.2f, want %.2f", result, tt.expected)
			}
		})
	}
}

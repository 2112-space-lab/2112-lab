package xpolygon

import "math"

type Edge struct {
	Start Point
	End   Point
}

type Point struct {
	Latitude  float64 // Latitude (Y-axis)
	Longitude float64 // Longitude (X-axis)
}

// Point3D represents a point in 3D space.
type Point3D struct {
	Latitude  float64
	Longitude float64
	Altitude  float64
}

// Edge3D represents an edge in 3D space.
type Edge3D struct {
	Start Point3D
	End   Point3D
}

// IsPointInPolygon checks if a given point (latitude, longitude) is inside a polygon
func IsPointInPolygon(point Point, polygon []Point) bool {
	// Ray-casting algorithm to determine if a point is inside a polygon
	// The algorithm counts how many times a ray starting from the point intersects the polygon's edges

	intersections := 0
	n := len(polygon)

	// Iterate over each edge of the polygon
	for i := 0; i < n; i++ {
		// Get the current edge
		p1 := polygon[i]
		p2 := polygon[(i+1)%n] // Loop back to the first point

		// Check if the point lies on the edge (ignoring vertical lines for simplicity)
		if point.Latitude > math.Min(p1.Latitude, p2.Latitude) && point.Latitude <= math.Max(p1.Latitude, p2.Latitude) {
			if point.Longitude <= math.Max(p1.Longitude, p2.Longitude) {
				if p1.Latitude != p2.Latitude {
					// Compute the x coordinate of the intersection
					xIntersection := (point.Latitude-p1.Latitude)*(p2.Longitude-p1.Longitude)/(p2.Latitude-p1.Latitude) + p1.Longitude
					if p1.Longitude == p2.Longitude || point.Longitude <= xIntersection {
						intersections++
					}
				}
			}
		}
	}

	// If the number of intersections is odd, the point is inside the polygon
	return intersections%2 == 1
}

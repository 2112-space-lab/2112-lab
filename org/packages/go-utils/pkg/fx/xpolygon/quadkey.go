package xpolygon

import (
	"fmt"
	"math"

	"github.com/org/2112-space-lab/org/go-utils/pkg/fx/xconstants"
)

type Quadkey struct {
	Latitude  float64
	Longitude float64
	Level     int
}

func NewQuadkey(lat float64, long float64, level int) Quadkey {
	return Quadkey{
		Latitude:  lat,
		Longitude: long,
		Level:     level,
	}
}

func (q *Quadkey) Key() string {
	return fmt.Sprintf("%d-%f-%f", q.Level, q.Latitude, q.Longitude)
}

// Constants

// DistanceTo computes the great-circle distance between two Quadkeys using the haversine formula
func (q *Quadkey) DistanceTo(other Point) float64 {
	// Convert latitudes and longitudes from degrees to radians
	lat1 := q.Latitude * xconstants.PI_DIVIDE_BY_180
	lon1 := q.Longitude * xconstants.PI_DIVIDE_BY_180
	lat2 := other.Latitude * xconstants.PI_DIVIDE_BY_180
	lon2 := other.Longitude * xconstants.PI_DIVIDE_BY_180

	// Compute differences
	dLat := lat2 - lat1
	dLon := lon2 - lon1

	// Haversine formula
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Distance in kilometers
	return xconstants.EARTH_RADIUS_KM * c
}

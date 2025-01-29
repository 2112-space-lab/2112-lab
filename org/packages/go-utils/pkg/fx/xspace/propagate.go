package xspace

import (
	"fmt"
	"math"
	"time"

	"github.com/joshuaferrara/go-satellite"
	xpolygon "github.com/org/2112-space-lab/org/go-utils/pkg/fx/xpolygon"
)

// PropagateSatellite propagates the satellite's position to the specified time.
// Returns a QuadKey for the satellite's position at the given time or an error.
func PropagateSatellite(tleLine1, tleLine2 string, t time.Time) (xpolygon.Quadkey, satellite.Satellite, error) {
	// Create satellite record from TLE lines
	satrec := satellite.TLEToSat(tleLine1, tleLine2, satellite.GravityWGS84)

	year, month, day := t.Date()
	hour, minute, second := t.Clock()

	position, _ := satellite.Propagate(satrec, year, int(month), day, hour, minute, second)
	if satrec.Error != 0 {
		return xpolygon.Quadkey{}, satellite.Satellite{}, fmt.Errorf("propagation error code: %d", satrec.Error)
	}

	// Calculate GST for ECI to LLA conversion
	gmst := satellite.GSTimeFromDate(year, int(month), day, hour, minute, second)

	// Convert ECI to Geodetic (lat, lon, alt)
	altitude, _, geoPosition := satellite.ECIToLLA(position, gmst)

	quadKey := xpolygon.NewQuadkey(geoPosition.Latitude, geoPosition.Longitude, int(altitude))
	return quadKey, satrec, nil
}

// SatellitePosition represents a satellite's position at a given time.
type SatellitePosition struct {
	Latitude  float64   // Degrees
	Longitude float64   // Degrees
	Altitude  float64   // Kilometers
	Time      time.Time // Timestamp
}

func PropagateRange(tleLine1, tleLine2 string, start, end time.Time, interval time.Duration) ([]SatellitePosition, error) {
	// Create satellite record from TLE lines
	satrec := satellite.TLEToSat(tleLine1, tleLine2, satellite.GravityWGS84)
	if satrec.Error != 0 {
		return nil, fmt.Errorf("TLE to Satellite error code: %d", satrec.Error)
	}

	var positions []SatellitePosition

	// Iterate through the time range
	for current := start; current.Before(end) || current.Equal(end); current = current.Add(interval) {
		// Extract date and time components for propagation
		year, month, day := current.Date()
		hour, minute, second := current.Clock()

		// Propagate the satellite's position
		position, _ := satellite.Propagate(satrec, year, int(month), day, hour, minute, second)
		if satrec.Error != 0 {
			return nil, fmt.Errorf("propagation error code: %d at %v", satrec.Error, current)
		}

		// Calculate GMST for ECI to LLA conversion
		gmst := satellite.GSTimeFromDate(year, int(month), day, hour, minute, second)

		// Convert ECI to Geodetic (lat, lon, alt)
		altitude, _, geoPosition := satellite.ECIToLLA(position, gmst)

		// Convert radians to degrees for latitude and longitude
		latitudeDeg := geoPosition.Latitude * (180.0 / math.Pi)
		longitudeDeg := geoPosition.Longitude * (180.0 / math.Pi)

		// Append the calculated position to the result
		positions = append(positions, SatellitePosition{
			Latitude:  latitudeDeg,
			Longitude: longitudeDeg,
			Altitude:  altitude,
			Time:      current,
		})
	}

	return positions, nil
}

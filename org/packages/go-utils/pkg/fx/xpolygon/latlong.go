package xpolygon

import "github.com/org/2112-space-lab/org/go-utils/pkg/fx/xconstants"

// Coordinate represents a geographic coordinate value
type Coordinate struct {
	C float64 // Coordinate value in degrees
}

// ToRadians converts the Coordinate object from degrees to radians
func (c Coordinate) ToRadians() float64 {
	return c.C * xconstants.PI_DIVIDE_BY_180
}

// ToDegrees converts the Coordinate object from radians to degrees
func (c Coordinate) ToDegrees() float64 {
	return c.C * xconstants.I180_DIVIDE_BY_PI
}

// LatLong represents a geographic coordinate pair in WGS84 format
type LatLong struct {
	Lat Coordinate // Latitude in degrees
	Lon Coordinate // Longitude in degrees
}

// LatRadians returns the latitude in radians
func (ll LatLong) LatRadians() float64 {
	return ll.Lat.ToRadians()
}

// LonRadians returns the longitude in radians
func (ll LatLong) LonRadians() float64 {
	return ll.Lon.ToRadians()
}

// LatDegrees returns the latitude in degrees
func (ll LatLong) LatDegrees() float64 {
	return ll.Lat.C
}

// LonDegrees returns the longitude in degrees
func (ll LatLong) LonDegrees() float64 {
	return ll.Lon.C
}

// ToRadians converts both latitude and longitude to radians
func (ll LatLong) ToRadians() LatLong {
	return LatLong{
		Lat: Coordinate{C: ll.Lat.ToRadians()},
		Lon: Coordinate{C: ll.Lon.ToRadians()},
	}
}

// ToDegrees converts both latitude and longitude to degrees
func (ll LatLong) ToDegrees() LatLong {
	return LatLong{
		Lat: Coordinate{C: ll.Lat.ToDegrees()},
		Lon: Coordinate{C: ll.Lon.ToDegrees()},
	}
}

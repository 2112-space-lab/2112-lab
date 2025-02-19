package xspace

import (
	"math"
	"time"

	xconstants "github.com/org/2112-space-lab/org/go-utils/pkg/fx/xconstants"
	xpolygon "github.com/org/2112-space-lab/org/go-utils/pkg/fx/xpolygon"
)

type OrbitType string

const (
	OrbitTypeUnknown OrbitType = "UNKNOWN"
	OrbitTypeLEO     OrbitType = "LEO" // Low Earth Orbit: 160 - 2000 km
	OrbitTypeMEO     OrbitType = "MEO" // Medium Earth Orbit: 2000 - 35786 km
	OrbitTypeGEO     OrbitType = "GEO" // Geostationary Orbit: ~35786 km
	OrbitTypeHEO     OrbitType = "HEO" // High Earth Orbit: > 35786 km
)

// LatLonToCartesian converts latitude, longitude, and altitude to Cartesian coordinates
func LatLonToCartesian(latitude, longitude, altitude float64) (float64, float64, float64) {

	// Convert latitude and longitude from degrees to radians
	latRad := DegreesToRadians(latitude)
	lonRad := DegreesToRadians(longitude)

	// Calculate the Cartesian coordinates using the spherical to Cartesian transformation
	x := (xconstants.EARTH_RADIUS_KM + altitude) * math.Cos(latRad) * math.Cos(lonRad)
	y := (xconstants.EARTH_RADIUS_KM + altitude) * math.Cos(latRad) * math.Sin(lonRad)
	z := (xconstants.EARTH_RADIUS_KM + altitude) * math.Sin(latRad)

	return x, y, z
}

// Normalize normalizes a vector
func Normalize(x, y, z float64) (float64, float64, float64) {
	magnitude := math.Sqrt(x*x + y*y + z*z)
	return x / magnitude, y / magnitude, z / magnitude
}

// DotProduct calculates the dot product of two vectors
func DotProduct(x1, y1, z1, x2, y2, z2 float64) float64 {
	return x1*x2 + y1*y2 + z1*z2
}

// CalculateIntegratedElevationFromPoint computes the elevation of a satellite relative to a ground point
func CalculateIntegratedElevationFromPoint(satellitePos xpolygon.Point, satelliteAltKm float64, groundPoint xpolygon.Point) float64 {
	// Convert the latitude/longitude of the ground point and satellite to 3D Cartesian coordinates
	groundX, groundY, groundZ := LatLonToCartesian(groundPoint.Latitude, groundPoint.Longitude, 0)
	satX, satY, satZ := LatLonToCartesian(satellitePos.Latitude, satellitePos.Longitude, satelliteAltKm)

	// Compute the vector from the ground to the satellite
	vecX, vecY, vecZ := satX-groundX, satY-groundY, satZ-groundZ

	// Compute the vector from the ground to Earth's center (opposite of the ground point)
	earthVecX, earthVecY, earthVecZ := Normalize(-groundX, -groundY, -groundZ)

	// Normalize the satellite vector
	vecX, vecY, vecZ = Normalize(vecX, vecY, vecZ)

	// Special case handling for direct overhead (tolerance for floating point precision)
	if math.Abs(vecX-earthVecX) < xconstants.EPSILON && math.Abs(vecY-earthVecY) < xconstants.EPSILON && math.Abs(vecZ-earthVecZ) < xconstants.EPSILON {
		return 90.0
	}

	// Compute the dot product of the vectors
	dotProd := DotProduct(vecX, vecY, vecZ, earthVecX, earthVecY, earthVecZ)

	// Ensure the cosine value is between -1 and 1 for valid arc cosine computation
	if dotProd > 1.0 {
		dotProd = 1.0
	} else if dotProd < -1.0 {
		dotProd = -1.0
	}

	// Compute the elevation angle in radians
	elevationRad := math.Acos(dotProd)

	// Convert the elevation angle from radians to degrees
	elevation := RadiansToDegrees(elevationRad)

	// Elevation cannot be more than 90 degrees
	if elevation > 90 {
		elevation = 90.0
	}

	return elevation
}

// Helper function to convert degrees to radians
func DegreesToRadians(degrees float64) float64 {
	return degrees * xconstants.PI_DIVIDE_BY_180
}

// Helper function to convert radians to degrees
func RadiansToDegrees(radians float64) float64 {
	return radians * 180.0 / math.Pi
}

// ComputeAverageAltitude computes avarage altitude from apogee and perigeee
func ComputeAverageAltitude(apogee, perigee float64) float64 {
	apogeeActual := apogee + xconstants.EARTH_RADIUS_KM
	perigeeActual := perigee + xconstants.EARTH_RADIUS_KM
	return (apogeeActual + perigeeActual) / 2
}

func calculateFraction(altitude float64) float64 {
	switch {
	case altitude < 200: // LEO
		return 0.01
	case altitude < 3578: // MEO
		return 0.05
	default: // GEO
		return 0.1
	}
}

// CalculateOptimalTimestep calculates optimal timesteps
func CalculateOptimalTimestep(altitude, tileRadius float64) time.Duration {

	orbitalVelocity := math.Sqrt(xconstants.GM / (xconstants.EARTH_RADIUS + altitude)) // Orbital velocity (m/s)
	timeOverTile := tileRadius / orbitalVelocity                                       // Time to cross a tile (seconds)
	fraction := calculateFraction(altitude)
	optimalTimestep := timeOverTile * fraction

	return time.Duration(optimalTimestep) * time.Second
}

// ComputeOrbitType determines the orbit type based on altitude in kilometers.
func ComputeOrbitType(alt float64) OrbitType {

	switch {
	case alt < 160:
		return OrbitTypeUnknown
	case alt < 2000:
		return OrbitTypeLEO
	case alt < 35786:
		return OrbitTypeMEO
	case alt == 35786:
		return OrbitTypeGEO
	default:
		return OrbitTypeHEO
	}
}

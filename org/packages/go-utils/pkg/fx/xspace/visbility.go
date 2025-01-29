package xspace

import (
	"log"
	"math"
	"time"

	"github.com/joshuaferrara/go-satellite"
	xconstants "github.com/org/2112-space-lab/org/go-utils/pkg/fx/xconstants"
	xpolygon "github.com/org/2112-space-lab/org/go-utils/pkg/fx/xpolygon"
)

// ComputeVisibilityWindow computes the visibility window for a satellite over a given tile.
func ComputeVisibilityWindow(
	noradID, tleLine1, tleLine2 string,
	point xpolygon.Point,
	radius float64,
	startTime, endTime time.Time, timeStep time.Duration,
) (time.Time, float64) {

	maxElevation := -1.0

	tileRadiusKm := radius / 1000

	satrec := satellite.TLEToSat(tleLine1, tleLine2, satellite.GravityWGS84)

	aos := ComputeAOS(satrec, point, tileRadiusKm, startTime, endTime, timeStep, &maxElevation)
	if aos.IsZero() {
		return aos, maxElevation
	} else {
		log.Printf("AOS detected at %v for NORAD ID %s", aos, noradID)
	}

	return aos, maxElevation
}

// ComputeAOS computes the Acquisition of Signal (AOS) time for a satellite over a given tile.
func ComputeAOS(
	satrec satellite.Satellite, point xpolygon.Point,
	tileRadiusKm float64, startTime, endTime time.Time,
	timeStep time.Duration, maxElevation *float64,
) time.Time {

	for t := startTime; t.Before(endTime); t = t.Add(timeStep) {
		altitude, geo, err := PropagateSatellitePosition(satrec, t)
		if err != nil {
			log.Printf("[ERROR] Propagation failed at %v: %v", t, err)
			continue
		}

		satellitePos := xpolygon.Point{Latitude: geo.Latitude, Longitude: geo.Longitude}

		if Intersects(point, satellitePos, tileRadiusKm, altitude) {

			elevation := CalculateIntegratedElevationFromPoint(satellitePos, altitude, point)

			// Check if AOS is valid
			if elevation > 0 {
				*maxElevation = math.Max(*maxElevation, elevation)
				return t
			}
		}
	}

	return time.Time{}
}

// ComputeLOS computes the Loss of Signal (LOS) time for a satellite over a given tile.
func ComputeLOS(
	satrec satellite.Satellite, point xpolygon.Point,
	tileRadiusKm float64, aos time.Time, endTime time.Time,
	timeStep time.Duration, maxElevation *float64,
) time.Time {

	for t := aos; t.Before(endTime); t = t.Add(timeStep) {
		altitude, geo, err := PropagateSatellitePosition(satrec, t)
		if err != nil {
			log.Printf("[ERROR] Propagation failed at %v: %v", t, err)
			continue
		}

		if altitude > 10000 {
			return aos
		}

		satellitePos := xpolygon.Point{Latitude: geo.Latitude, Longitude: geo.Longitude}

		if !Intersects(point, satellitePos, tileRadiusKm, altitude) {
			return t
		}
	}

	return time.Time{}
}

func Intersects(tileCenter xpolygon.Point, satellitePos xpolygon.Point, tileRadiusKm float64, altitude float64) bool {
	// Compute the distance from the satellite to the tile center point (ignoring altitude)
	centerDistance := HaversineDistance(satellitePos.Latitude, satellitePos.Longitude, tileCenter.Latitude, tileCenter.Longitude, 0, 0)

	// Add a small margin of error to account for floating-point precision issues
	marginOfError := 0.01

	// Check if the satellite is within the tile radius from the center of the tile
	return centerDistance <= tileRadiusKm+marginOfError
}

// PropagateSatellitePosition calculates the satellite's geodetic position at a specific time.
func PropagateSatellitePosition(satrec satellite.Satellite, t time.Time) (float64, satellite.LatLong, error) {
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	position, _ := satellite.Propagate(satrec, year, int(month), day, hour, minute, second)
	gmst := satellite.GSTimeFromDate(year, int(month), day, hour, minute, second)
	altitude, _, geo := satellite.ECIToLLA(position, gmst)
	return altitude, geo, nil
}

// CalculateIntegratedElevation computes the elevation of a satellite relative to a ground point.
func CalculateIntegratedElevation(satelliteQuadKey xpolygon.Quadkey, satelliteAltitude float64, pointQuadKey xpolygon.Point) float64 {
	return CalculateIntegratedElevationFromPoint(
		xpolygon.Point{Latitude: satelliteQuadKey.Latitude, Longitude: satelliteQuadKey.Longitude},
		satelliteAltitude,
		xpolygon.Point{Latitude: pointQuadKey.Latitude, Longitude: pointQuadKey.Longitude},
	)
}

func HaversineDistance(lat1, lon1, lat2, lon2, altitude1, altitude2 float64) float64 {
	dLat := DegreesToRadians(lat2 - lat1)
	dLon := DegreesToRadians(lon2 - lon1)
	lat1Rad := DegreesToRadians(lat1)
	lat2Rad := DegreesToRadians(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	surfaceDistance := xconstants.EARTH_RADIUS_KM * c

	// Adjust for altitude (adding z-axis distance)
	altitudeDiff := altitude2 - altitude1
	return math.Sqrt(surfaceDistance*surfaceDistance + altitudeDiff*altitudeDiff)
}

// Compute the satellite's visible region (horizon) at a given time.
func ComputeSatelliteHorizon(t time.Time, line1 string, line2 string) ([]xpolygon.Point, error) {
	// Propagate satellite position to get the current location
	satrec := satellite.TLEToSat(line1, line2, satellite.GravityWGS84)
	altitude, geo, err := PropagateSatellitePosition(satrec, t)
	if err != nil {
		return nil, err
	}

	// Sub-satellite point (directly beneath the satellite)
	subSatellitePoint := xpolygon.Point{Latitude: geo.Latitude, Longitude: geo.Longitude}

	// Calculate the horizon distance (approximated by the satellite's altitude and Earth's radius)
	horizonDistance := math.Sqrt(2 * xconstants.EARTH_RADIUS_KM * altitude)

	// Now, we will define the visible region as a circle with the calculated horizon distance
	// The circular region is represented by multiple points (approximated here as a set of points around the horizon)
	visibleRegion := make([]xpolygon.Point, 0)

	// Number of points to represent the circle (for simplicity, let's use 36 points for a full circle)
	numPoints := 36
	for i := 0; i < numPoints; i++ {
		angle := float64(i) * (2 * math.Pi / float64(numPoints))
		latOffset := horizonDistance / xconstants.EARTH_RADIUS_KM * math.Sin(angle)
		lonOffset := horizonDistance / xconstants.EARTH_RADIUS_KM * math.Cos(angle)

		// Calculate the new latitude and longitude by applying the offsets
		lat := subSatellitePoint.Latitude + latOffset
		lon := subSatellitePoint.Longitude + lonOffset

		// Add the point to the visible region
		visibleRegion = append(visibleRegion, xpolygon.Point{Latitude: lat, Longitude: lon})
	}

	// Return the visible region as a xpolygon (simplified)
	return visibleRegion, nil
}

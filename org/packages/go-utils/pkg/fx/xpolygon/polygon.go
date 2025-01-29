package xpolygon

import (
	"log"
	"math"

	"github.com/org/2112-space-lab/org/go-utils/pkg/fx/xconstants"
)

// Polygon represents a geographic polygon.

type Polygon struct {
	Center     Quadkey // The center Quadkey of the polygon
	NbFaces    int     // Number of faces in the polygon
	Radius     float64 // Radius of the polygon in meters
	Boundaries []Point // Vertices (points) of the polygon's boundary
}

// NewPolygon creates a new Polygon given the number of faces, center coordinates, zoom level, and radius.
func NewPolygon(nbFaces int, center LatLong, level int, radius float64) Polygon {
	boundaries := generateBoundaries(nbFaces, center, radius)
	return Polygon{
		NbFaces:    nbFaces,
		Radius:     radius,
		Center:     NewQuadkey(center.LatDegrees(), center.LonDegrees(), level),
		Boundaries: boundaries,
	}
}

// generateBoundaries calculates the boundary points of the polygon.
func generateBoundaries(nbFaces int, center LatLong, radius float64) []Point {
	boundaries := make([]Point, 0, nbFaces)

	centerLatRad := center.LatRadians()
	centerLonRad := center.LonRadians()

	angleStep := 2 * math.Pi / float64(nbFaces)

	for i := 0; i < nbFaces; i++ {
		angle := angleStep * float64(i)

		deltaLat := radius * math.Cos(angle) / xconstants.EARTH_RADIUS
		deltaLon := radius * math.Sin(angle) / (xconstants.EARTH_RADIUS * math.Cos(centerLatRad))

		newLat := centerLatRad + deltaLat
		newLon := centerLonRad + deltaLon

		newLatDeg := newLat * xconstants.I180_DIVIDE_BY_PI
		newLonDeg := newLon * xconstants.I180_DIVIDE_BY_PI

		boundaries = append(boundaries, Point{
			Latitude:  newLatDeg,
			Longitude: newLonDeg,
		})
	}

	return boundaries
}

// calculateZoomLevelForTileRadius determines the optimal zoom level for a given tile radius.
func calculateZoomLevelForTileRadius(tileRadius float64) int {
	for zoom := 0; zoom <= xconstants.MAX_ZOOM_LEVEL; zoom++ {
		tileSize := xconstants.EARTH_CIRCUMFERENCE_METER / math.Pow(2, float64(zoom))

		log.Printf("Zoom: %d, Tile Size: %.2f, Half Tile Size: %.2f", zoom, tileSize, tileSize/2)
		if tileSize/2 <= tileRadius {
			return zoom
		}
	}
	return xconstants.MAX_ZOOM_LEVEL
}

// calculateTileRadiusForZoom calculates the exact radius of a tile given a specific zoom level.
func calculateTileRadiusForZoom(zoom int) float64 {
	tileSize := xconstants.EARTH_CIRCUMFERENCE_METER / math.Pow(2, float64(zoom))
	tileRadius := tileSize / 2
	return tileRadius
}

// GenerateAllTilesForRadius generates all tiles and their polygons for a given radius.
func GenerateAllTilesForRadius(tileRadius float64, nbFaces int) map[Quadkey]Polygon {
	// Calculate the zoom level based on the tile radius
	zoom := calculateZoomLevelForTileRadius(tileRadius)
	numTiles := int(math.Pow(2, float64(zoom))) // Total number of tiles at that zoom level
	radius := calculateTileRadiusForZoom(zoom)

	// Determine the range for X and Y coordinates of tiles
	startX := 0
	endX := numTiles

	// Map to hold the generated tile polygons
	tilePolygons := make(map[Quadkey]Polygon)

	// Iterate over all tile X and Y coordinates at the given zoom level
	for x := startX; x < endX; x++ {
		for y := 0; y < numTiles; y++ {
			// Convert tile (x, y) coordinates to latitude and longitude
			lat, lon := TileXYToLatLon(x, y, zoom)

			// Create the center coordinates of the tile
			center := LatLong{Lat: Coordinate{lat}, Lon: Coordinate{lon}}

			// Generate the quadkey for the center of the tile
			centerQuadkey := NewQuadkey(center.LatDegrees(), center.LonDegrees(), zoom)

			// Create the polygon for the tile using the given radius and number of faces
			polygon := NewPolygon(nbFaces, center, zoom, radius)

			// Store the polygon with its corresponding quadkey
			tilePolygons[centerQuadkey] = polygon
		}
	}

	// Return the generated polygons
	return tilePolygons
}

// TileXYToLatLon converts tile coordinates to lat/lon.
func TileXYToLatLon(x, y, zoom int) (float64, float64) {
	n := math.Pow(2, float64(zoom))
	lon := float64(x)/n*360.0 - 180.0

	latRad := math.Atan(math.Sinh(math.Pi * (1 - 2*float64(y)/n)))
	lat := latRad * xconstants.I180_DIVIDE_BY_PI

	return lat, lon
}

// LatLonToTileXY converts lat/lon to tile coordinates.
func LatLonToTileXY(lat, lon float64, zoom int) (int, int) {
	n := math.Pow(2, float64(zoom))

	tileX := int((lon + 180.0) / 360.0 * n)
	tileY := int((1.0 - math.Log(math.Tan(lat*math.Pi/180.0)+1.0/math.Cos(lat*math.Pi/180.0))/math.Pi) / 2.0 * n)

	return tileX, tileY
}

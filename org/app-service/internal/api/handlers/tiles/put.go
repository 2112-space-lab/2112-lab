package tiles

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"

	api_mappers "github.com/org/2112-space-lab/org/app-service/pkg/api"
	xconstants "github.com/org/2112-space-lab/org/go-utils/pkg/fx/xconstants"
)

const TILE_URL = "https://%s.basemaps.cartocdn.com/light_all/%d/%d/%d.png" // CartoDB Positron URL

// FetchTile fetches a tile for the given zoom level and geographical latitude/longitude.
func FetchTile(zoom int, lat float64, lon float64) (*api_mappers.RawTile, error) {
	// Ensure the zoom level is valid
	if zoom <= 0 {
		return nil, fmt.Errorf("invalid zoom level")
	}

	// Convert latitude and longitude to tile coordinates (x, y)
	x, y := latLonToTileCoordinates(lat, lon, zoom)

	// Choose a random subdomain for load balancing (a, b, or c)
	subdomains := []string{"a", "b", "c"}
	subdomain := subdomains[rand.Intn(len(subdomains))]

	// Construct the URL for fetching the tile using fmt.Sprintf
	url := fmt.Sprintf(TILE_URL, subdomain, zoom, x, y)

	// Make the HTTP request to fetch the tile
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tile: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response is successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch tile: HTTP status %d", resp.StatusCode)
	}

	// Read the tile data (this is either PNG for raster tiles or PBF for vector tiles)
	tileData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read tile data: %v", err)
	}

	// Return the RawTile object with calculated x, y coordinates
	return &api_mappers.RawTile{
		ZoomLevel: zoom,
		X:         x,
		Y:         y,
		TileData:  tileData,
	}, nil
}

// latLonToTileCoordinates converts latitude and longitude to tile coordinates (x, y) at a given zoom level
func latLonToTileCoordinates(lat float64, lon float64, zoom int) (int, int) {
	// Number of tiles at the zoom level (2^zoom)
	n := math.Pow(2, float64(zoom))

	// Calculate tile X and Y coordinates from longitude and latitude
	x := int((lon + 180.0) / 360.0 * n)
	latRad := lat * xconstants.PI_DIVIDE_BY_180
	y := int((1.0 - math.Log(math.Tan(latRad)+1/math.Cos(latRad))/math.Pi) / 2.0 * n)

	return x, y
}

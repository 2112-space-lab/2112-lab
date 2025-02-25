package api_mappers

// RawTile represents the raw tile data fetched from a tile server
type RawTile struct {
	ZoomLevel int
	X         int
	Y         int
	TileData  []byte // Raw data for the tile (e.g., PBF for vector tiles, PNG for raster tiles)
}

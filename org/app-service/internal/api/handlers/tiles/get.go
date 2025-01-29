package tiles

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/domain"
	"github.com/Elbujito/2112/src/app-service/internal/services"
	"github.com/labstack/echo/v4"
	xconstants "github.com/org/2112-space-lab/org/go-utils/pkg/fx/xconstants"
)

type TileHandler struct {
	Service services.TileService
}

// NewTileHandler creates a new handler with the provided TileService.
func NewTileHandler(service services.TileService) *TileHandler {
	return &TileHandler{Service: service}
}

// GetAllTiles fetches all available tiles.
func (h *TileHandler) GetAllTiles(c echo.Context) error {
	tiles, err := h.Service.FindAllTiles(c.Request().Context(), "todoTileHandler")
	if err != nil {
		c.Echo().Logger.Error("Failed to fetch tiles: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to fetch tiles")
	}

	// If no tiles are found
	if len(tiles) == 0 {
		c.Echo().Logger.Error(xconstants.ERROR_ID_NOT_FOUND)
		return xconstants.ERROR_ID_NOT_FOUND
	}

	// Return tiles in the response
	return c.JSON(http.StatusOK, tiles)
}

// GetTilesInRegionHandler handles requests to fetch tiles in a region.
func (h *TileHandler) GetTilesInRegionHandler(c echo.Context) error {
	// Parse query parameters for bounding box
	minLatStr := c.QueryParam("minLat")
	minLonStr := c.QueryParam("minLon")
	maxLatStr := c.QueryParam("maxLat")
	maxLonStr := c.QueryParam("maxLon")

	// Convert query parameters to float64
	minLat, err := strconv.ParseFloat(minLatStr, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid minLat parameter")
	}
	minLon, err := strconv.ParseFloat(minLonStr, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid minLon parameter")
	}
	maxLat, err := strconv.ParseFloat(maxLatStr, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid maxLat parameter")
	}
	maxLon, err := strconv.ParseFloat(maxLonStr, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid maxLon parameter")
	}

	// Call the service to fetch tiles
	tiles, err := h.Service.GetTilesInRegion(c.Request().Context(), "todoTileHandler", minLat, minLon, maxLat, maxLon)
	if err != nil {
		c.Logger().Error("Failed to fetch tiles in region:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to fetch tiles in region")
	}

	// Return tiles in JSON response
	return c.JSON(http.StatusOK, tiles)
}

// GetPaginatedSatelliteMappings fetches a paginated list of satellite mappings with optional search filters.
func (h *TileHandler) GetPaginatedSatelliteMappings(c echo.Context) error {
	// Parse query parameters for pagination
	pageStr := c.QueryParam("page")
	pageSizeStr := c.QueryParam("pageSize")
	searchWildcard := c.QueryParam("search") // Retrieve optional search query

	// Convert parameters to integers with defaults
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1 // Default to page 1 if invalid
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10 // Default to 10 records per page if invalid
	}

	// Create SearchRequest object
	searchRequest := &domain.SearchRequest{
		Wildcard: searchWildcard,
	}

	// Call the service method for pagination with search filters
	mappings, totalRecords, err := h.Service.ListSatellitesMappingWithPagination(c.Request().Context(), "todoTileHandler", page, pageSize, searchRequest)
	if err != nil {
		c.Echo().Logger.Error("Failed to fetch paginated satellites mappings: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to fetch satellites mappings")
	}

	// Prepare the response
	response := map[string]interface{}{
		"totalRecords": totalRecords,
		"page":         page,
		"pageSize":     pageSize,
		"mappings":     mappings,
	}

	return c.JSON(http.StatusOK, response)
}

// GetSatelliteMappingByNoradID handles requests to fetch tiles in a region.
func (h *TileHandler) GetSatelliteMappingsByNoradID(c echo.Context) error {
	// Parse query parameters for bounding box
	noradID := c.QueryParam("noradID")

	// Call the service to fetch mappings
	mappings, err := h.Service.GetSatelliteMappingsByNoradID(c.Request().Context(), "todoTileHandler", noradID)
	if err != nil {
		c.Logger().Error("Failed to fetch mappings:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to fetch mappings by norad ID [%s]", noradID)
	}

	// Return tiles in JSON response
	return c.JSON(http.StatusOK, mappings)
}

// RecomputeMappingsByNoradID handles requests to recompute satellite mappings for a given NORAD ID.
func (h *TileHandler) RecomputeMappingsByNoradID(c echo.Context) error {
	// Extract the NORAD ID from the query parameter
	noradID := c.QueryParam("noradID")
	if noradID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing noradID parameter")
	}

	// Extract startTime and endTime from query parameters
	startTimeStr := c.QueryParam("startTime")
	endTimeStr := c.QueryParam("endTime")

	if startTimeStr == "" || endTimeStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Both startTime and endTime parameters are required")
	}

	// Parse the startTime and endTime parameters
	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid startTime format, expected RFC3339")
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid endTime format, expected RFC3339")
	}

	// Ensure startTime is before endTime
	if !startTime.Before(endTime) {
		return echo.NewHTTPError(http.StatusBadRequest, "startTime must be before endTime")
	}

	// Call the service method to recompute mappings
	err = h.Service.RecomputeMappings(c.Request().Context(), "todoTileHandler", noradID, startTime, endTime)
	if err != nil {
		c.Logger().Error("Failed to recompute mappings for NORAD ID:", noradID, "Error:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to recompute mappings for NORAD ID")
	}

	// Return a success response
	return c.JSON(http.StatusOK, map[string]string{
		"message":   "Mappings recomputed successfully",
		"noradID":   noradID,
		"startTime": startTime.Format(time.RFC3339),
		"endTime":   endTime.Format(time.RFC3339),
	})
}

package tiles

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/org/2112-space-lab/org/app-service/internal/config/constants"
	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	"github.com/org/2112-space-lab/org/app-service/internal/services"
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
		c.Echo().Logger.Error(constants.ERROR_ID_NOT_FOUND)
		return constants.ERROR_ID_NOT_FOUND
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

// GetSatelliteMappingsBySpaceID handles requests to fetch tiles in a region.
func (h *TileHandler) GetSatelliteMappingsBySpaceID(c echo.Context) error {
	// Parse query parameters for bounding box
	spaceID := c.QueryParam("spaceID")

	// Call the service to fetch mappings
	mappings, err := h.Service.GetSatelliteMappingsBySpaceID(c.Request().Context(), "todoTileHandler", spaceID)
	if err != nil {
		c.Logger().Error("Failed to fetch mappings:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to fetch mappings by space ID [%s]", spaceID)
	}

	// Return tiles in JSON response
	return c.JSON(http.StatusOK, mappings)
}

// RecomputeMappingsBySpaceID handles requests to recompute satellite mappings for a given SPACE ID.
func (h *TileHandler) RecomputeMappingsBySpaceID(c echo.Context) error {
	// Extract the SPACE ID from the query parameter
	spaceID := c.QueryParam("spaceID")
	if spaceID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing spaceID parameter")
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
	err = h.Service.RecomputeMappings(c.Request().Context(), "todoTileHandler", spaceID, startTime, endTime)
	if err != nil {
		c.Logger().Error("Failed to recompute mappings for SPACE ID:", spaceID, "Error:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to recompute mappings for SPACE ID")
	}

	// Return a success response
	return c.JSON(http.StatusOK, map[string]string{
		"message":   "Mappings recomputed successfully",
		"spaceID":   spaceID,
		"startTime": startTime.Format(time.RFC3339),
		"endTime":   endTime.Format(time.RFC3339),
	})
}

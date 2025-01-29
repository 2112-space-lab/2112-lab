package apiaudittrail

import (
	"net/http"
	"strconv"

	"github.com/Elbujito/2112/src/app-service/internal/services"
	"github.com/labstack/echo/v4"
)

// AuditTrailHandler handles API requests related to AuditTrails.
type AuditTrailHandler struct {
	Service services.AuditTrailService
}

// NewAuditTrailHandler creates a new handler with the provided AuditTrailService.
func NewAuditTrailHandler(service services.AuditTrailService) *AuditTrailHandler {
	return &AuditTrailHandler{Service: service}
}

// GetAuditTrails retrieves audit trails by record ID and table name or with pagination.
func (h *AuditTrailHandler) GetAuditTrails(c echo.Context) error {
	// Optional query parameters
	tableName := c.QueryParam("tableName")
	recordID := c.QueryParam("recordID")

	// Pagination parameters
	pageStr := c.QueryParam("page")
	pageSizeStr := c.QueryParam("pageSize")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// If tableName and recordID are provided, retrieve specific audit trails
	if tableName != "" && recordID != "" {
		trails, err := h.Service.GetByRecordIDAndTable(c.Request().Context(), tableName, recordID)
		if err != nil {
			c.Echo().Logger.Error("Failed to retrieve audit trails: ", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Unable to fetch audit trails")
		}
		return c.JSON(http.StatusOK, trails)
	}

	// Otherwise, retrieve paginated audit trails
	trails, total, err := h.Service.GetAllWithPagination(c.Request().Context(), page, pageSize)
	if err != nil {
		c.Echo().Logger.Error("Failed to fetch paginated audit trails: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to fetch audit trails")
	}

	// Prepare the response
	response := map[string]interface{}{
		"page":        page,
		"pageSize":    pageSize,
		"totalCount":  total,
		"auditTrails": trails,
	}

	return c.JSON(http.StatusOK, response)
}

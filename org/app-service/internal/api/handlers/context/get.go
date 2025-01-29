package apicontext

import (
	"net/http"
	"strconv"

	"github.com/Elbujito/2112/src/app-service/internal/domain"
	"github.com/Elbujito/2112/src/app-service/internal/services"
	"github.com/labstack/echo/v4"
)

// ContextHandler handles API requests related to GameContexts.
type ContextHandler struct {
	Service services.ContextService
}

// NewContextHandler creates a new handler with the provided ContextService.
func NewContextHandler(service services.ContextService) *ContextHandler {
	return &ContextHandler{Service: service}
}

// CreateContext handles the creation of a new GameContext.
func (h *ContextHandler) CreateContext(c echo.Context) error {
	var gameContext domain.GameContext

	// Bind JSON request body to the domain.GameContext struct
	if err := c.Bind(&gameContext); err != nil {
		c.Echo().Logger.Error("Failed to bind GameContext: ", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	// Call the service to create the context
	createdContext, err := h.Service.Create(c.Request().Context(), gameContext)
	if err != nil {
		c.Echo().Logger.Error("Failed to create GameContext: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to create context")
	}

	return c.JSON(http.StatusCreated, createdContext)
}

// UpdateContext handles updating an existing GameContext.
func (h *ContextHandler) UpdateContext(c echo.Context) error {
	var gameContext domain.GameContext

	// Bind JSON request body to the domain.GameContext struct
	if err := c.Bind(&gameContext); err != nil {
		c.Echo().Logger.Error("Failed to bind GameContext: ", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	// Call the service to update the context
	updatedContext, err := h.Service.Update(c.Request().Context(), gameContext)
	if err != nil {
		c.Echo().Logger.Error("Failed to update GameContext: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to update context")
	}

	return c.JSON(http.StatusOK, updatedContext)
}

// GetContextByName retrieves a GameContext by its unique name.
func (h *ContextHandler) GetContextByName(c echo.Context) error {
	name := c.Param("name") // Extract context name from the URL path

	// Call the service to retrieve the context
	context, err := h.Service.GetByUniqueName(c.Request().Context(), domain.GameContextName(name))
	if err != nil {
		c.Echo().Logger.Error("Failed to retrieve GameContext: ", err)
		return echo.NewHTTPError(http.StatusNotFound, "Context not found")
	}

	return c.JSON(http.StatusOK, context)
}

// DeleteContextByName deletes a GameContext by its unique name.
func (h *ContextHandler) DeleteContextByName(c echo.Context) error {
	name := c.Param("name") // Extract context name from the URL path

	// Call the service to delete the context
	if err := h.Service.DeleteByUniqueName(c.Request().Context(), domain.GameContextName(name)); err != nil {
		c.Echo().Logger.Error("Failed to delete GameContext: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to delete context")
	}

	return c.NoContent(http.StatusNoContent)
}

// ActivateContext activates a GameContext by its unique name.
func (h *ContextHandler) ActivateContext(c echo.Context) error {
	name := c.Param("name") // Extract context name from the URL path

	// Call the service to activate the context
	if err := h.Service.ActiveContext(c.Request().Context(), domain.GameContextName(name)); err != nil {
		c.Echo().Logger.Error("Failed to activate GameContext: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to activate context")
	}

	return c.NoContent(http.StatusNoContent)
}

// DeactivateContext deactivates a GameContext by its unique name.
func (h *ContextHandler) DeactivateContext(c echo.Context) error {
	name := c.Param("name") // Extract context name from the URL path

	// Call the service to deactivate the context
	if err := h.Service.DisableContext(c.Request().Context(), domain.GameContextName(name)); err != nil {
		c.Echo().Logger.Error("Failed to deactivate GameContext: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to deactivate context")
	}

	return c.NoContent(http.StatusNoContent)
}

// GetPaginatedContexts retrieves paginated GameContexts with optional search.
func (h *ContextHandler) GetPaginatedContexts(c echo.Context) error {
	// Parse query parameters
	pageStr := c.QueryParam("page")
	pageSizeStr := c.QueryParam("pageSize")
	wildcard := c.QueryParam("search") // Optional search parameter

	// Convert parameters to integers with defaults
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// Retrieve paginated contexts
	contexts, err := h.Service.FindAllWithPagination(c.Request().Context(), page, pageSize, wildcard)
	if err != nil {
		c.Echo().Logger.Error("Failed to fetch paginated contexts: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to fetch contexts")
	}

	// Prepare response
	response := map[string]interface{}{
		"page":       page,
		"pageSize":   pageSize,
		"contexts":   contexts,
		"totalCount": len(contexts), // Adjust if total count is separately available
	}

	return c.JSON(http.StatusOK, response)
}

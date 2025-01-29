package apiuser

import (
	"net/http"
	"strconv"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/labstack/echo/v4"
)

// UserHandler handles API requests related to users.
type UserHandler struct{}

// NewUserHandler creates a new handler for user-related operations.
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetUsers retrieves a paginated list of users from Clerk.
func (h *UserHandler) GetUsers(c echo.Context) error {
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

	offset := (page - 1) * pageSize

	listParams := &user.ListParams{
		ListParams: clerk.ListParams{
			Limit:  clerk.Int64(int64(pageSize)),
			Offset: clerk.Int64(int64(offset)),
		},
	}

	users, err := user.List(c.Request().Context(), listParams)
	if err != nil {
		c.Echo().Logger.Error("Failed to retrieve users from Clerk: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to fetch users")
	}

	// Prepare the response
	response := map[string]interface{}{
		"page":       page,
		"pageSize":   pageSize,
		"totalCount": users.TotalCount,
		"users":      users.Users,
	}

	return c.JSON(http.StatusOK, response)
}

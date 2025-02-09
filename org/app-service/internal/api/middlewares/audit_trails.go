package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/labstack/echo/v4"
	"github.com/org/2112-space-lab/org/app-service/internal/services"
)

// GenerateRouteTableMapping dynamically generates a route-to-table mapping from Echo routes.
func GenerateRouteTableMapping(e *echo.Echo) map[string]string {
	mapping := make(map[string]string)

	// Iterate over all registered routes
	for _, route := range e.Routes() {
		// Derive the table name from the route's path
		// Example: "/satellites/orbit" -> "satellites"
		path := route.Path
		segments := strings.Split(strings.Trim(path, "/"), "/")
		if len(segments) > 0 {
			tableName := segments[0] // Use the first segment as the table name
			mapping[path] = tableName
		}
	}

	return mapping
}

// LogNonGETRequestsMiddleware logs all non-GET HTTP requests to the audit trail using AuditTrailService.
func LogNonGETRequestsMiddleware(routeTableMapping map[string]string, auditTrailService services.AuditTrailService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method != echo.GET {
				req := c.Request()

				claims, ok := clerk.SessionClaimsFromContext(req.Context())
				userID := "unknown"
				if ok && claims != nil {
					userID = claims.Subject // The Subject contains the Clerk user ID
				}

				// Determine the table name based on the request path
				tableName, found := routeTableMapping[c.Path()]
				if !found {
					tableName = "unknown"
				}

				// Read and capture request body
				var requestBody map[string]interface{}
				if req.ContentLength > 0 && (req.Method == echo.POST || req.Method == echo.PUT || req.Method == echo.DELETE) {
					bodyBytes, err := io.ReadAll(req.Body)
					if err != nil {
						c.Logger().Errorf("Failed to read request body: %v", err)
					} else {
						// Restore the body to the request after reading
						req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
						_ = json.Unmarshal(bodyBytes, &requestBody)
					}
				}

				// Use LogAudit method to log the request
				err := auditTrailService.LogAudit(
					req.Context(),
					tableName,            // Table name determined from path
					c.Request().URL.Path, // Record ID (use the request path)
					c.Request().Method,   // Action (HTTP method)
					userID,               // Performed by (user ID from Clerk)
					requestBody,          // Changes (parsed request body)
				)

				if err != nil {
					c.Logger().Errorf("Failed to log audit trail: %v", err)
				}
			}

			return next(c)
		}
	}
}

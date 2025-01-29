package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SlashesMiddleware returns Slashes Middleware
func SlashesMiddleware() echo.MiddlewareFunc {
	return middleware.RemoveTrailingSlash()
}

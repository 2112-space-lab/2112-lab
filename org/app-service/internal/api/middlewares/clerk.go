package middlewares

import (
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/labstack/echo/v4"
)

// ClerkMiddleware integrates Clerk authentication into Echo
func ClerkMiddleware(opts ...clerkhttp.AuthorizationOption) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			handler := clerkhttp.RequireHeaderAuthorization(opts...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.SetRequest(r)
				c.Response().Writer = w

				claims, ok := clerk.SessionClaimsFromContext(r.Context())
				if ok {
					c.Set("claims", claims)
				}

				_ = next(c)
			}))

			handler.ServeHTTP(c.Response().Writer, c.Request())
			return nil
		}
	}
}

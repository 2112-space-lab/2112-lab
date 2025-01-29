package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/org/2112-space-lab/org/go-utils/pkg/fx/xutils"
)

// BodyDumpMiddleware returns Body dump Middleware
func BodyDumpMiddleware() echo.MiddlewareFunc {
	return middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		if len(reqBody) > 0 {
			obj, err := xutils.PrettyJSONString(string(reqBody))
			if err != nil {
				c.Logger().Error("Error unmarshalling request body: ", err)
				return
			}
			c.Logger().Debug("Request Body: ", obj)
		}
	})
}

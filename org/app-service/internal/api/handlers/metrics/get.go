package metrics

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func GetMetrics(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/plain; charset=utf-8")
	promhttp.Handler().ServeHTTP(c.Response(), c.Request())
	return nil
}

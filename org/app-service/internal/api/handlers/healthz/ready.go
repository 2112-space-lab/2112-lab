package healthz

import (
	"net/http"

	"github.com/Elbujito/2112/src/app-service/internal/api/handlers"
	"github.com/Elbujito/2112/src/app-service/internal/api/helpers"
	"github.com/Elbujito/2112/src/app-service/internal/clients/dbc"

	"github.com/labstack/echo/v4"
)

func Ready(c echo.Context) error {
	dbClient := dbc.GetDBClient()
	if err := dbClient.Ping(); err != nil {
		return helpers.Error(c, err, nil)
	}

	payload := map[string]string{
		"message": "ready",
	}
	return c.JSON(http.StatusOK, handlers.Success(payload))
}

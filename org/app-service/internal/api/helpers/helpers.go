package helpers

import (
	"github.com/labstack/echo/v4"
)

func init() {
}

func Error(c echo.Context, err error, ori_err error) error {
	c.Logger().Error(err, ori_err)
	return err
}

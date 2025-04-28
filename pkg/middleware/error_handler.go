package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/pkg/response"
)

func GlobalErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	if he, ok := err.(*echo.HTTPError); ok {
		c.JSON(he.Code, map[string]interface{}{
			"success": false,
			"error": map[string]interface{}{
				"message": he.Message,
			},
		})
		return
	}

	response.Error(c, err)
}

package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/config/logger"
	"github.com/yusuffugurlu/go-project/pkg/errors"
)

type StandardResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func Success(c echo.Context, status int, data interface{}) error {
	if status == 0 {
		status = http.StatusOK
	}

	return c.JSON(status, StandardResponse{
		Success: true,
		Data:    data,
	})
}

func Error(c echo.Context, err error) error {
	var statusCode int
	var errorResponse map[string]interface{}

	if appErr, ok := errors.AsAppError(err); ok {
		statusCode = appErr.Code
		errorResponse = map[string]interface{}{
			"message": appErr.Message,
		}

		if appErr.Details != nil {
			errorResponse["details"] = appErr.Details
		}

		if statusCode >= 500 {
			logger.Log.Errorf("Internal server error: %s", appErr.Message)
		}
	} else {
		statusCode = http.StatusInternalServerError
		errorResponse = map[string]interface{}{
			"message": "Internal server error",
		}
		logger.Log.Errorf("Unhandled error: %v", err)
	}

	return c.JSON(statusCode, StandardResponse{
		Success: false,
		Error:   errorResponse,
	})
}

func Created(c echo.Context, data interface{}) error {
	return Success(c, http.StatusCreated, data)
}

func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/internal/services"
	"github.com/yusuffugurlu/go-project/pkg/response"
)

type LogController interface {
	GetAllLogs(e echo.Context) error
}

type logController struct {
	auditLogService services.AuditLogService
}

func NewLogController(auditLogService services.AuditLogService) LogController {
	return &logController{auditLogService: auditLogService}
}

func (l *logController) GetAllLogs(e echo.Context) error {
	logs, err := l.auditLogService.GetAllAuditLogs()
	if err != nil {
		return err
	}

	return response.Success(e, http.StatusOK, logs)
}
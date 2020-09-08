package api

import (
	"net/http"

	"github.com/chadhao/logit/modules/log/service"
	"github.com/labstack/echo/v4"
)

// queryLogs 获取logs
func queryLogs(c echo.Context) error {
	req := new(service.QueryLogsInput)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	resp, err := service.QueryLogs(req)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, resp.Logs)
}

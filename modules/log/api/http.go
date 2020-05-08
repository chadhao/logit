package api

import (
	"net/http"

	"github.com/chadhao/logit/modules/log/model"
	"github.com/labstack/echo/v4"
)

// queryLogs 获取logs
func queryLogs(c echo.Context) error {
	req := new(model.QueryLog)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	logs, err := req.Find()
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, logs)
}

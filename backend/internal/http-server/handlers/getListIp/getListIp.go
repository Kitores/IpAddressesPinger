package getListIp

import (
	"IpAddressPinger/backend/internal/database/postgresSql"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

type ListIpGetter interface {
	GetListIp() ([]postgresSql.Data, error)
}

func New(log *slog.Logger, getter ListIpGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		functionName := "http-server/handlers/getListIp.New"
		log = log.With(slog.String("funcName", functionName))
		log.Debug("start")
		listIp, err := getter.GetListIp()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, listIp)
	}
}

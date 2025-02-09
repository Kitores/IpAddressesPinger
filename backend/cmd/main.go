package main

import (
	"IpAddressPinger/backend/internal/config"
	"IpAddressPinger/backend/internal/database/postgresSql"
	"IpAddressPinger/backend/internal/http-server/handlers/getListIp"
	"IpAddressPinger/backend/internal/http-server/handlers/postPingInfo"
	"IpAddressPinger/backend/internal/setupLogger"
	"IpAddressPinger/backend/lib/logger/sl"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger.SetupLogger(cfg.Env)

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.UserDb, cfg.Password, cfg.Dbname, cfg.SSLmode)
	storage, err := postgresSql.NewPg(connStr)
	if err != nil || storage == nil {
		log.Error("Failed to initialize storage: %v", sl.Err(err))
	}

	e := echo.New()

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Разрешаем запросы с любых источников
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
	})
	e.GET("/getListIp", getListIp.New(log, storage))
	e.POST("/pingInfo", postPingInfo.New(log, storage))
	e.Start(cfg.Address)
}

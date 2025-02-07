package main

import (
	"IpAddressPinger/backend/internal/database/postgresSql"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func main() {

	conn := ""
	storage, err := postgresSql.NewPg(conn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(storage)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}

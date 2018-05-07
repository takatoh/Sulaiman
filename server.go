package main

import (
//	"net/http"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.File("static/html/index.html")
	})

	e.Static("/css", "static/css")

	e.Logger.Fatal(e.Start(":1323"))
}

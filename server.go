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
	e.Static("/js", "static/js")
	e.Static("/img", "photos/img")
	e.Static("/thumb", "photos/thumb")

	e.GET("/list/:page", func(c echo.Context) error {
		jsonFile := "list" + c.Param("page") + ".json"
		return c.File(jsonFile)
	})

	e.Logger.Fatal(e.Start(":1323"))
}

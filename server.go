package main

import (
//	"net/http"

	"github.com/labstack/echo"

	"github.com/takatoh/sulaiman/handler"
)

func main() {
	e := echo.New()

	e.GET("/", handler.IndexGet)

	e.Static("/css", "static/css")
	e.Static("/js", "static/js")
	e.Static("/img", "photos/img")
	e.Static("/thumb", "photos/thumb")

	e.GET("/list/:page", handler.ListGet)

	e.Logger.Fatal(e.Start(":1323"))
}

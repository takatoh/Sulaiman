package handler

import (
	"github.com/labstack/echo"
)

func IndexGet(c echo.Context) error {
	return c.File("static/html/index.html")
}

func ListGet(c echo.Context) error {
	jsonFile := "list" + c.Param("page") + ".json"
		return c.File(jsonFile)
}

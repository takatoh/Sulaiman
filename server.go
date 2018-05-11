package main

import (
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/takatoh/sulaiman/handler"
	"github.com/takatoh/sulaiman/data"
)

func main() {
	db, err := gorm.Open("sqlite3", "sulaiman.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&data.Photo{})

	e := echo.New()
	h := handler.New(db)

	e.GET("/", h.IndexGet)
	e.GET("/title", h.TitleGet)

	e.Static("/css", "static/css")
	e.Static("/js", "static/js")
	e.Static("/img", "photos/img")
	e.Static("/thumb", "photos/thumb")

	e.GET("/list/:page", h.ListGet)
	e.POST("/upload", h.UploadPost)

	e.Logger.Fatal(e.Start(":1323"))
}

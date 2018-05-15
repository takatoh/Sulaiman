package main

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"flag"

	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/takatoh/sulaiman/handler"
	"github.com/takatoh/sulaiman/data"
)

func main() {
	opt_config := flag.String("c", "config.json", "Specify config file.")
	flag.Parse()

	var config = new(data.Config)
	jsonString, err := ioutil.ReadFile(*opt_config)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(jsonString, config)

	db, err := gorm.Open("sqlite3", "sulaiman.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&data.Photo{})

	e := echo.New()
	h := handler.New(db, config)

	e.GET("/", h.Index)
	e.GET("/title", h.Title)

	e.Static("/css", "static/css")
	e.Static("/js", "static/js")
	e.Static("/img", config.PhotoDir + "/img")
	e.Static("/thumb", config.PhotoDir + "/thumb")

	e.GET("/list/:page", h.List)
	e.POST("/upload", h.Upload)

	port := ":" + strconv.Itoa(config.Port)
	e.Logger.Fatal(e.Start(port))
}

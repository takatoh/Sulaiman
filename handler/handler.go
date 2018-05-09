package handler

import (
	"os"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

func IndexGet(c echo.Context) error {
	return c.File("static/html/index.html")
}

func ListGet(c echo.Context) error {
	jsonFile := "list" + c.Param("page") + ".json"
		return c.File(jsonFile)
}

func UploadPost(c echo.Context) error {
	file, _ := c.FormFile("file")
	src, _ := file.Open()
	defer src.Close()

	dst, _ := os.Create("photos/img/img101.jpg")
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	res := UploadResponse{
		Status: "OK",
		Photo: Photo{
			ID: 101,
			Url: "http://localhost:1323/img/img101.jpg",
			Thumb: "http://localhost:1323/thumb/thumb101.jpg",
		},
	}

	return c.JSON(http.StatusOK, res)
}

type UploadResponse struct {
	Status string `json:"status"`
	Photo  Photo  `json:"photo"`
}

type Photo struct {
	ID    int    `json:"id"`
	Url   string `json:"url"`
	Thumb string `json:"thumb"`
}
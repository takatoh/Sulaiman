package handler

import (
	"os"
	"io"
	"net/http"
	"image"
	"image/jpeg"
	_ "image/png"
	_ "image/gif"

	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	"github.com/nfnt/resize"
)

type Handler struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Handler {
	p := new(Handler)
	p.db = db
	return p
}

func (h *Handler) IndexGet(c echo.Context) error {
	return c.File("static/html/index.html")
}

func (h *Handler) ListGet(c echo.Context) error {
	jsonFile := "list" + c.Param("page") + ".json"
		return c.File(jsonFile)
}

func (h *Handler) UploadPost(c echo.Context) error {
	file, _ := c.FormFile("file")
	src, _ := file.Open()
	defer src.Close()

	dst, _ := os.Create("photos/img/img101.jpg")
	defer dst.Close()

	_, _ = io.Copy(dst, src)

	img := "img/img101.jpg"
	thumb := makeThumbnail(img)

	res := UploadResponse{
		Status: "OK",
		Photo: Photo{
			ID: 101,
			Url: "http://localhost:1323/" + img,
			Thumb: "http://localhost:1323/" + thumb,
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

func makeThumbnail(src_file string) string {
	src, _ := os.Open("photos/" + src_file)
	defer src.Close()

	config, _, _ := image.DecodeConfig(src)
	src.Seek(0, 0)
	img, _, _ := image.Decode(src)

	var resized_img image.Image
	if config.Width >= config.Height {
		resized_img = resize.Resize(120, 0, img, resize.Lanczos3)
	} else {
		resized_img = resize.Resize(0, 120, img, resize.Lanczos3)
	}
	thumb, _ := os.Create("photos/thumb/thumb101.jpg")
	jpeg.Encode(thumb, resized_img, nil)
	thumb.Close()

	return "thumb/thumb101.jpg"
}

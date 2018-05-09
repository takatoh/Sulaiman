package handler

import (
	"os"
	"io"
	"net/http"
	"image"
	"image/jpeg"
	_ "image/png"
	_ "image/gif"
	"strconv"

	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	"github.com/nfnt/resize"

	"github.com/takatoh/sulaiman/data"
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
	page, _ := strconv.Atoi(c.Param("page"))
	offset := page - 1
	var photos []data.Photo
	h.db.Order("id desc").Offset(offset).Limit(10).Find(&photos)

	var resPhotos []*Photo
	for _, p := range photos {
		resPhotos = append(resPhotos, newPhoto(p.ID, p.ImagePath, p.ThumbPath))
	}
	res := ListResponse{
		Status: "OK",
		Page: page,
		Next: "/list/" + strconv.Itoa(page + 1),
		Photos: resPhotos,
	}

	return c.JSON(http.StatusOK, res)
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

type ListResponse struct {
	Status string   `json:"status"`
	Page   int      `json:"page"`
	Next   string   `json:"next"`
	Photos []*Photo `json:"photos"`
}

type UploadResponse struct {
	Status string `json:"status"`
	Photo  Photo  `json:"photo"`
}

type Photo struct {
	ID    uint   `json:"id"`
	Url   string `json:"url"`
	Thumb string `json:"thumb"`
}

func newPhoto(id uint, img, thumb string) *Photo {
	p := new(Photo)
	p.ID = id
	p.Url = "http://localhost:1323/" + img
	p.Thumb = "http://localhost:1323/" + thumb
	return p
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

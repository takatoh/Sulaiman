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
	"strings"

	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	"github.com/nfnt/resize"

	"github.com/takatoh/sulaiman/data"
)

type Handler struct {
	db     *gorm.DB
	config *data.Config
}

func New(db *gorm.DB, config *data.Config) *Handler {
	p := new(Handler)
	p.db = db
	p.config = config
	return p
}

func (h *Handler) Index(c echo.Context) error {
	return c.File("static/html/index.html")
}

func (h *Handler) Title(c echo.Context) error {
	return c.String(http.StatusOK, h.config.SiteName)
}

func (h *Handler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.Param("page"))
	offset := (page - 1) * 10
	var photos []data.Photo
	h.db.Order("id desc").Offset(offset).Limit(10).Find(&photos)

	var resPhotos []*ResPhoto
	for _, p := range photos {
		resPhotos = append(
			resPhotos,
			newResPhoto(
				p.ID,
				buildURL(p.ImagePath, h.config),
				"/" + p.ImagePath,
				"/" + p.ThumbPath,
			),
		)
	}
	var next string
	if len(resPhotos) < 10 {
		next = ""
	} else {
		next = "/list/" + strconv.Itoa(page + 1)
	}
	res := ListResponse{
		Status: "OK",
		Page: page,
		Next: next,
		Photos: resPhotos,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) Upload(c echo.Context) error {
	file, _ := c.FormFile("file")
	filename := file.Filename
	src, _ := file.Open()
	defer src.Close()

	var lastPhoto data.Photo
	h.db.Last(&lastPhoto)
	newId := int(lastPhoto.ID) + 1
	pos := strings.LastIndex(filename, ".")
	ext := filename[pos:]
	img := "img/img" + strconv.Itoa(newId) + ext

	dst, _ := os.Create(h.config.PhotoDir + "/" + img)
	defer dst.Close()

	_, _ = io.Copy(dst, src)

	thumb := makeThumbnail(h.config.PhotoDir, img, newId)

	deleteKey := c.FormValue("key")
	newPhoto := data.Photo{ ImagePath: img, ThumbPath: thumb, DeleteKey: deleteKey }
	h.db.Create(&newPhoto)

	res := UploadResponse{
		Status: "OK",
		Photo: newResPhoto(
			newPhoto.ID,
			buildURL(img, h.config),
			"/" + img,
			"/" + thumb,
		),
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) Delete(c echo.Context) error {
	var photo data.Photo
	id, _ := strconv.Atoi(c.FormValue("id"))
	deleteID := uint(id)
	deleteKey := c.FormValue("key")
	photo.ID = deleteID
	h.db.First(&photo)
	var res DeleteResponse
	if photo.DeleteKey == deleteKey {
		deletePhoto(photo, h.config)
		h.db.Delete(&photo)
		res = DeleteResponse{
			Status: "OK",
			PhotoID: deleteID,
		}
	} else {
		res = DeleteResponse{
			Status: "NG",
			PhotoID: deleteID,
		}
	}

	return c.JSON(http.StatusOK, res)
}

type ListResponse struct {
	Status string      `json:"status"`
	Page   int         `json:"page"`
	Next   string      `json:"next"`
	Photos []*ResPhoto `json:"photos"`
}

type UploadResponse struct {
	Status string     `json:"status"`
	Photo  *ResPhoto  `json:"photo"`
}

type ResPhoto struct {
	ID    uint   `json:"id"`
	Url   string `json:"url"`
	Img   string `json:"img"`
	Thumb string `json:"thumb"`
}

func newResPhoto(id uint, url, img, thumb string) *ResPhoto {
	p := new(ResPhoto)
	p.ID = id
	p.Url = url
	p.Img = img
	p.Thumb = thumb
	return p
}

type DeleteResponse struct {
	Status  string `json:"status"`
	PhotoID uint   `json:"photo_id"`
}

func makeThumbnail(photo_dir, src_file string, id int) string {
	src, _ := os.Open(photo_dir + "/" + src_file)
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
	thumb_file := "thumb/thumb" + strconv.Itoa(id) + ".jpg"
	thumb, _ := os.Create(photo_dir + "/" + thumb_file)
	jpeg.Encode(thumb, resized_img, nil)
	thumb.Close()

	return thumb_file
}

func buildURL(path string, config *data.Config) string {
	return "http://" + config.HostName + "/" + path
}

func deletePhoto(photo data.Photo, config *data.Config) {
	photo_dir := config.PhotoDir
	os.Remove(photo_dir + "/" + photo.ImagePath)
	os.Remove(photo_dir + "/" + photo.ThumbPath)
	return
}
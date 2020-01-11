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
	"time"
	"fmt"

	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	"github.com/nfnt/resize"

	"github.com/takatoh/sulaiman/data"
)

const (
	photosPerPage = 25
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
	offset := (page - 1) * photosPerPage
	var photos []data.Photo
	h.db.Order("id desc").Offset(offset).Limit(photosPerPage).Find(&photos)

	var resPhotos []*ResPhoto
	for _, p := range photos {
		resPhotos = append(
			resPhotos,
			newResPhoto(
				p.ID,
				buildURL(p.ImagePath, h.config),
				"/" + p.ImagePath,
				"/" + p.ThumbPath,
				p.Width,
				p.Height,
				p.FileSize,
				p.CreatedAt,
			),
		)
	}
	var next string
	if len(resPhotos) < photosPerPage {
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

	photo := data.Photo{}
	h.db.Create(&photo)
	newId := int(photo.ID)
	pos := strings.LastIndex(filename, ".")
	ext := filename[pos:]
	img := "img/img" + strconv.Itoa(newId) + ext

	dst, _ := os.Create(h.config.PhotoDir + "/" + img)
	defer dst.Close()

	_, _ = io.Copy(dst, src)

	f, _ := dst.Stat()
	filesize := f.Size()
	thumb, width, height := makeThumbnail(h.config.PhotoDir, img, newId)

	deleteKey := c.FormValue("key")
	photo.ImagePath = img
	photo.ThumbPath = thumb
	photo.DeleteKey = deleteKey
	photo.Width     = width
	photo.Height    = height
	photo.FileSize  = filesize
	h.db.Save(&photo)

	var photos []data.Photo
	var count int
	var deletePhotoID uint
	h.db.Find(&photos).Count(&count)
	if h.config.MaxPhotoCount > 0 {
		if count > h.config.MaxPhotoCount {
			var first data.Photo
			h.db.First(&first)
			deletePhotoID = first.ID
			deletePhoto(first, h.config)
			h.db.Delete(&first)
		}
	}

	res := UploadResponse{
		Status: "OK",
		Photo: newResPhoto(
			photo.ID,
			buildURL(img, h.config),
			"/" + img,
			"/" + thumb,
			photo.Width,
			photo.Height,
			photo.FileSize,
			photo.CreatedAt,
		),
		DeletePhotoID: deletePhotoID,
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

func (h *Handler) Count(c echo.Context) error {
	var photos []data.Photo
	var count int
	h.db.Find(&photos).Count(&count)
	res := CountResponse{Count: count}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) First(c echo.Context) error {
	var photo data.Photo
	h.db.First(&photo)
	res := newResPhoto(
		photo.ID,
		buildURL(photo.ImagePath, h.config),
		"/" + photo.ImagePath,
		"/" + photo.ThumbPath,
		photo.Width,
		photo.Height,
		photo.FileSize,
		photo.CreatedAt,
	)

	return c.JSON(http.StatusOK, res)
}

type ListResponse struct {
	Status string      `json:"status"`
	Page   int         `json:"page"`
	Next   string      `json:"next"`
	Photos []*ResPhoto `json:"photos"`
}

type UploadResponse struct {
	Status        string    `json:"status"`
	Photo         *ResPhoto `json:"photo"`
	DeletePhotoID uint      `json:"deletePhotoId"`
}

type ResPhoto struct {
	ID       uint   `json:"id"`
	Url      string `json:"url"`
	Img      string `json:"img"`
	Thumb    string `json:"thumb"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize string `json:"filesize"`
	Posted   string `json:"posted"`
}

func newResPhoto(id uint, url, img, thumb string, width, height int, filesize int64, posted time.Time) *ResPhoto {
	p := new(ResPhoto)
	p.ID = id
	p.Url = url
	p.Img = img
	p.Thumb = thumb
	p.Width = width
	p.Height = height
	p.FileSize = humanBytes(filesize)
	p.Posted = posted.Format("2006-01-02 15:04:05 -07:00")
	return p
}

type DeleteResponse struct {
	Status  string `json:"status"`
	PhotoID uint   `json:"photoId"`
}

type CountResponse struct {
	Count int `json:"count"`
}

func makeThumbnail(photo_dir, src_file string, id int) (string, int, int) {
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

	return thumb_file, config.Width, config.Height
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

func humanBytes(filesize int64) string {
	size := float64(filesize)
	units := []string{ "bytes", "KB", "MB" }
	u := 0
	for size >= 1000.0 {
		u = u + 1
		size = size / 1000.0
	}

	var hsize string
	if size > 100.0 {
		hsize = fmt.Sprintf("%.0f %s", size, units[u])
	} else if size > 10.0 {
		hsize = fmt.Sprintf("%.1f %s", size, units[u])
	} else {
		hsize = fmt.Sprintf("%.2f %s", size, units[u])
	}
	return hsize
}

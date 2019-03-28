package data

import (
	"github.com/jinzhu/gorm"
)

type Photo struct {
	gorm.Model
	ImagePath string
	ThumbPath string
	DeleteKey string
	Width     int
	Heihgt    int
	FileSize  int
}

type Config struct {
	SiteName      string `json:"site_name"`
	HostName      string `json:"host_name"`
	Port          int    `json:"port"`
	PhotoDir      string `json:"photo_dir"`
	MaxPhotoCount int    `json:"max_photo_count"`
}

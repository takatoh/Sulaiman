package data

import (
	"github.com/jinzhu/gorm"
)

type Photo struct {
	gorm.Model
	ImagePath string
	ThumbPath string
}

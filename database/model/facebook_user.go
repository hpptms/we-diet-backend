package model

import "gorm.io/gorm"

type FacebookUser struct {
	gorm.Model
	FacebookID      string `gorm:"uniqueIndex"`
	FacebookEmail   string
	FacebookName    string
	FacebookPicture string
}

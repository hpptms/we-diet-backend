package database

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName   string `gorm:"uniqueIndex"`
	Password   string
	Subscribe  bool
	Permission int

	// Google認証用
	GoogleID string `gorm:"uniqueIndex"`
	Email    string
	Name     string
	Picture  string

	// Facebook認証用
	FacebookID      string `gorm:"uniqueIndex"`
	FacebookEmail   string
	FacebookName    string
	FacebookPicture string

	// TikTok認証用
	TikTokID     string `gorm:"uniqueIndex"`
	TikTokName   string
	TikTokAvatar string
}

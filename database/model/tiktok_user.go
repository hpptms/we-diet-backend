package model

import "gorm.io/gorm"

type TikTokUser struct {
	gorm.Model
	TikTokID     string `gorm:"uniqueIndex"`
	TikTokName   string
	TikTokAvatar string
}

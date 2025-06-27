package model

type TikTokUser struct {
	ID           uint   `gorm:"primaryKey"`
	TikTokID     string `gorm:"uniqueIndex"`
	TikTokName   string
	TikTokAvatar string

	Timestamps
}

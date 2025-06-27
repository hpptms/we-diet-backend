package model

type FacebookUser struct {
	ID              uint   `gorm:"primaryKey"`
	FacebookID      string `gorm:"uniqueIndex"`
	FacebookEmail   string
	FacebookName    string
	FacebookPicture string

	Timestamps
}

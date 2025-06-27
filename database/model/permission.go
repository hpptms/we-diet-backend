package model

type Permission struct {
	ID         uint   `gorm:"primaryKey"`
	Permission string `gorm:"uniqueIndex"`

	Timestamps
}

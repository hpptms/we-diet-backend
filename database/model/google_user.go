package model

type GoogleUser struct {
	ID       uint   `gorm:"primaryKey"`
	GoogleID string `gorm:"uniqueIndex"`
	Email    string
	Name     string
	Picture  string

	Timestamps
}

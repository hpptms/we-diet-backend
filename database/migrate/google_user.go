package migrate

import "gorm.io/gorm"

type GoogleUser struct {
	gorm.Model
	GoogleID string `gorm:"uniqueIndex"`
	Email    string
	Name     string
	Picture  string
}

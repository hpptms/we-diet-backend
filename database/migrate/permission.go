package migrate

import (
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Permission string `gorm:"uniqueIndex"`
}

package model

import "gorm.io/gorm"

type OtherService struct {
	gorm.Model
	ServiceName string `gorm:"uniqueIndex"`
	ServiceID   int    `gorm:"uniqueIndex"`
}

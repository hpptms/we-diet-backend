package migrate

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID              int `gorm:"uniqueIndex"`
	Email           *string
	UserName        string
	Password        *string
	Icon            *string
	Subscribe       bool `gorm:"default:0"`
	Permission      int  `gorm:"default:0"`
	OtherServices   int  `gorm:"default:0"` //1=google 2=facebook 3=tiktok
	OtherServicesId *string
}

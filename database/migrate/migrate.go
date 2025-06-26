package migrate

import (
	"my-gin-app/database"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&database.User{}, &Permission{})
}

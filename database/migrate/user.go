package migrate

import (
	"my-gin-app/database/model"

	"gorm.io/gorm"
)

// Userテーブルのマイグレーション
func MigrateUser(db *gorm.DB) error {
	return db.AutoMigrate(&model.User{})
}

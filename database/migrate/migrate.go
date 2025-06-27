package migrate

import (
	database "my-gin-app/database/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&database.User{},
		&database.FacebookUser{},
		&database.GoogleUser{},
		&database.TikTokUser{},
		&database.Permission{},
		&database.OtherService{},
	)
}

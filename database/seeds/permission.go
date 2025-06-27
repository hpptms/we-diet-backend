package seeds

import (
	"log"

	"my-gin-app/database/migrate"

	"gorm.io/gorm"
)

func PermissionSeed(db *gorm.DB) error {
	permissions := []struct {
		ID         int
		Permission string
	}{
		{555, "admin"},
		{0, "user"},
	}

	for _, perm := range permissions {
		var existing migrate.Permission
		err := db.Where("id = ?", perm.ID).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&migrate.Permission{Model: gorm.Model{ID: uint(perm.ID)}, Permission: perm.Permission}).Error; err != nil {
				log.Printf("Failed to seed permission '%s': %v", perm.Permission, err)
				return err
			} else {
				log.Printf("Seeded permission: %s (id=%d)", perm.Permission, perm.ID)
			}
		}
	}
	return nil
}

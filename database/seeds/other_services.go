package seeds

import (
	"log"

	"my-gin-app/database"
	"my-gin-app/database/migrate"

	"gorm.io/gorm"
)

func SeedOtherServices() {
	db := database.GetDB()
	services := []struct {
		Name string
		ID   int
	}{
		{"Google", 1},
		{"Facebook", 2},
		{"TikTok", 3},
	}

	for _, s := range services {
		var existing migrate.OtherService
		err := db.Where("service_id = ?", s.ID).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&migrate.OtherService{ServiceName: s.Name, ServiceID: s.ID}).Error; err != nil {
				log.Printf("Failed to seed other_service '%s': %v", s.Name, err)
			} else {
				log.Printf("Seeded other_service: %s", s.Name)
			}
		}
	}
}

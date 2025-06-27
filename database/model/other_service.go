package model

type OtherService struct {
	ID          uint   `gorm:"primaryKey"`
	ServiceName string `gorm:"uniqueIndex"`
	ServiceID   int    `gorm:"uniqueIndex"`

	Timestamps
}

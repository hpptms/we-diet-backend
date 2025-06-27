package model

type User struct {
	ID          uint   `gorm:"primaryKey"`
	UserName    string `gorm:"uniqueIndex"`
	Password    string
	Subscribe   bool
	Permission  int
	Picture     string
	ServiceName string
	ServiceID   int `gorm:"uniqueIndex"`

	Timestamps
}

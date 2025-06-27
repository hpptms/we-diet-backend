package main

import (
	"fmt"
	"os"

	migrate "my-gin-app/database/migrate"
	seeds "my-gin-app/database/seeds"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("DB接続失敗: " + err.Error())
	}

	if err := migrate.Migrate(db); err != nil {
		panic("DBマイグレーション失敗: " + err.Error())
	}
	fmt.Println("DBマイグレーション完了")

	// ユーザーseed投入
	if err := seeds.UserSeed(db); err != nil {
		panic("ユーザーseed投入失敗: " + err.Error())
	}
	fmt.Println("ユーザーseed投入完了")

	// パーミッションseed投入
	if err := seeds.PermissionSeed(db); err != nil {
		panic("パーミッションseed投入失敗: " + err.Error())
	}
	fmt.Println("パーミッションseed投入完了")

	// OtherService seed投入
	if err := seeds.OtherServiceSeed(db); err != nil {
		panic("OtherService seed投入失敗: " + err.Error())
	}
	fmt.Println("OtherService seed投入完了")
}

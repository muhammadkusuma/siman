package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Pastikan username:password dan nama db sesuai konfigurasi Anda
	dsn := "root:12345678@tcp(127.0.0.1:3306)/db_siman?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrasi semua model yang ada
	database.AutoMigrate(
		&User{},
		&Faculty{},
		&Department{},
		&Building{},
		&Room{},
		&AssetCategory{},
		&Asset{},
		&MaintenanceLog{},
		&MutationLog{},
		&AuditLog{},
	)

	DB = database
}
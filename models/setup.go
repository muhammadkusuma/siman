package models

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Pastikan username:password dan nama db sesuai konfigurasi Anda
	// Tambahkan parseTime=True agar struct time.Time bekerja
	dsn := "root:12345678@tcp(127.0.0.1:3306)/db_siman?charset=utf8mb4&parseTime=True&loc=Local"
	
	// Gunakan Logger agar error SQL terlihat di terminal
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	
	if err != nil {
		panic("Gagal koneksi ke database!")
	}

	// Migrasi satu per satu atau sekaligus.
	// PENTING: Asset harus sebelum MutationLog
	err = database.AutoMigrate(
		&User{},
		&Faculty{},
		&Department{},
		&Building{},
		&Room{},
		&AssetCategory{},
		&Asset{},
		&MaintenanceLog{},
		&MutationLog{}, // Ini yang sering bermasalah, pastikan tabel di atasnya sudah sukses
		&AuditLog{},
	)

	if err != nil {
		// Ini akan mencetak error spesifik MySQL jika migrasi gagal
		log.Fatal("GAGAL MIGRASI DATABASE: ", err)
	}

	log.Println("Database berhasil dimigrasi!")
	DB = database
}
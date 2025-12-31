package models

import (
	"log"

	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Ambil konfigurasi dari Environment Variables
	// Pastikan nama variabel ini SAMA PERSIS dengan yang ada di file .env
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Cek apakah variabel env terbaca (opsional, untuk debugging awal)
	if dbUser == "" || dbHost == "" {
		log.Println("Peringatan: Variabel lingkungan DB kosong. Pastikan file .env ada dan terbaca.")
	}

	// Buat DSN (Data Source Name) secara dinamis menggunakan fmt.Sprintf
	// Format: username:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// Gunakan Logger agar error SQL terlihat di terminal
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		// Log fatal akan menghentikan aplikasi jika koneksi gagal
		log.Fatal("GAGAL KONEKSI KE DATABASE! Cek konfigurasi .env anda. Error: ", err)
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

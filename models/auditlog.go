package models

import (
	"time"
)

// AuditLog mencatat aktivitas sensitif user (Siapa melakukan apa).
type AuditLog struct {
	ID     uint `json:"id" gorm:"primaryKey;autoIncrement;comment:ID Log Audit"`
	UserID uint `json:"user_id" gorm:"not null;comment:ID User pelaku aksi"`

	// Target Aksi
	Action    string `json:"action" gorm:"size:50;comment:Jenis Aksi (cth: CREATE, UPDATE, DELETE, LOGIN, MUTATION)"`
	TableName string `json:"table_name" gorm:"size:100;comment:Tabel yang dimodifikasi (cth: assets, users)"`
	RecordID  uint   `json:"record_id" gorm:"comment:ID Primary Key dari data yang diubah"`

	// Data Perubahan
	Changes string `json:"changes" gorm:"type:text;comment:JSON snapshot data (cth: {'old': 'Baik', 'new': 'Rusak'})"`

	IPAddress string `json:"ip_address" gorm:"size:50;comment:IP Address user saat kejadian (cth: 192.168.1.50)"`
	UserAgent string `json:"user_agent" gorm:"size:255;comment:Browser/Device yang digunakan user (cth: Chrome on Windows)"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;comment:Waktu kejadian"`

	// Relasi
	User User `json:"user" gorm:"foreignKey:UserID"`
}

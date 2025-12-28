package models

import (
	"time"
)

type AuditLog struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint      `json:"user_id" gorm:"not null;comment:Siapa yang melakukan aksi"`
	
	// Target Aksi
	Action    string    `json:"action" gorm:"size:50;comment:CREATE, UPDATE, DELETE, LOGIN, MUTATION"`
	TableName string    `json:"table_name" gorm:"size:100;comment:Tabel yang diubah (misal: assets)"`
	RecordID  uint      `json:"record_id" gorm:"comment:ID dari data yang diubah"`

	// Data Perubahan (Disimpan dalam format JSON String)
	// Contoh: {"old": "Baik", "new": "Rusak Berat"}
	Changes   string    `json:"changes" gorm:"type:text;comment:JSON snapshot data sebelum dan sesudah"`

	IPAddress string    `json:"ip_address" gorm:"size:50;comment:IP Address user saat akses"`
	UserAgent string    `json:"user_agent" gorm:"size:255;comment:Browser/Device yang digunakan"`
	
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relasi ke User
	User User `json:"user" gorm:"foreignKey:UserID"`
}
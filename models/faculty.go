package models

import (
	"time"
)

// Faculty merepresentasikan induk unit: Fakultas, Sekolah Pascasarjana, atau Direktorat.
type Faculty struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement;comment:ID Unik Fakultas/Unit Induk"`
	Code      string    `json:"code" gorm:"size:50;unique;comment:Kode Singkatan Unit (cth: FT, SPS, REK)"`
	Name      string    `json:"name" gorm:"size:255;not null;comment:Nama Lengkap Unit (cth: Fakultas Sains dan Teknologi)"`
	Type      string    `json:"type" gorm:"type:enum('Fakultas','Sekolah Pascasarjana','Direktorat','Lembaga');not null;comment:Jenis Institusi (cth: Fakultas)"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:Waktu pembuatan data"`
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:Waktu update data terakhir"`

	// Relasi
	Departments []Department `json:"departments" gorm:"foreignKey:FacultyID"`
}

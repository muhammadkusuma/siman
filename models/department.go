package models

import (
	"time"
)

// Department merepresentasikan Prodi (S1/S2/S3) atau Bagian/Biro spesifik.
type Department struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	FacultyID uint   `json:"faculty_id" gorm:"not null;comment:Foreign Key ke tabel Faculties"`
	
	Code      string `json:"code" gorm:"size:50;comment:Kode Prodi/Unit (misal: IF-S1)"`
	Name      string `json:"name" gorm:"size:255;not null;comment:Nama Prodi (misal: S1 Teknik Informatika)"`
	
	// Enum StudyLevel penting untuk membedakan aset Pascasarjana vs Sarjana
	StudyLevel string `json:"study_level" gorm:"type:enum('D3','S1','S2','S3','Profesi','Non-Akademik');default:'Non-Akademik';comment:Jenjang Studi. Pilih Non-Akademik untuk kantor admin."`
	
	HeadOfDepartment string `json:"head_of_department" gorm:"size:255;comment:Nama Kaprodi atau Kepala Unit"`

	// Relasi Belongs To
	Faculty Faculty `json:"faculty" gorm:"foreignKey:FacultyID"`
}
package models

import (
	"time"
)

// User menyimpan data pengguna aplikasi.
type User struct {
	ID           uint   `json:"id" gorm:"primaryKey;autoIncrement;comment:ID Unik User"`
	Username     string `json:"username" gorm:"size:100;unique;not null;comment:Username untuk login (cth: admin_tif)"`
	PasswordHash string `json:"-" gorm:"size:255;not null;comment:Password terenkripsi Bcrypt (Jangan diedit manual)"`
	FullName     string `json:"full_name" gorm:"size:255;comment:Nama Lengkap User (cth: Budi Santoso, S.Kom)"`
	Email        string `json:"email" gorm:"size:100;unique;comment:Email aktif untuk reset password (cth: budi@uin.ac.id)"`

	// Role Management
	Role string `json:"role" gorm:"type:enum('SuperAdmin','AdminFakultas','AdminProdi','Auditor','Peminjam');default:'Peminjam';comment:Hak Akses User (cth: AdminProdi)"`

	// Relasi Unit Kerja
	DepartmentID *uint `json:"department_id" gorm:"comment:Unit kerja user. Jika NULL berarti user tingkat Universitas/Pusat."`

	CreatedAt time.Time `json:"created_at" gorm:"comment:Waktu pendaftaran user"`
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:Waktu profil diupdate"`

	// Relasi
	Department *Department `json:"department" gorm:"foreignKey:DepartmentID"`
}

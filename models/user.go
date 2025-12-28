package models

import (
	"time"
)

type User struct {
	ID           uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string `json:"username" gorm:"size:100;unique;not null;comment:Username untuk login"`
	PasswordHash string `json:"-" gorm:"size:255;not null;comment:Password terenkripsi (Bcrypt)"` // JSON "-" agar tidak ikut ter-render di API
	FullName     string `json:"full_name" gorm:"size:255;comment:Nama Lengkap User"`
	Email        string `json:"email" gorm:"size:100;unique;comment:Email untuk reset password"`

	// Role Management
	Role string `json:"role" gorm:"type:enum('SuperAdmin','AdminFakultas','AdminProdi','Auditor','Peminjam');default:'Peminjam';comment:Hak Akses User"`

	// Relasi: User ini bertugas di Unit mana?
	// Menggunakan Pointer (*uint) agar bisa NULL (Contoh: SuperAdmin atau Auditor Pusat tidak terikat prodi)
	DepartmentID *uint `json:"department_id" gorm:"comment:Jika diisi, user hanya bisa kelola aset prodi ini"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relasi
	Department *Department `json:"department" gorm:"foreignKey:DepartmentID"`
}
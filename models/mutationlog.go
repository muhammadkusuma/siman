package models

import (
	"time"
)


// MutationLog mencatat sejarah perpindahan aset antar unit/ruangan.
type MutationLog struct {
	ID      uint `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetID uint `json:"asset_id" gorm:"not null"`

	// Relasi Asal (Pointer *uint karena bisa NULL jika barang baru input)
	FromDepartmentID *uint `json:"from_department_id" gorm:"comment:Unit Asal"`
	FromRoomID       *uint `json:"from_room_id" gorm:"comment:Ruangan Asal"`

	// Relasi Tujuan (Tidak boleh NULL)
	ToDepartmentID uint `json:"to_department_id" gorm:"not null;comment:Unit Tujuan"`
	ToRoomID       uint `json:"to_room_id" gorm:"comment:Ruangan Tujuan"`

	MutationDate time.Time `json:"mutation_date" gorm:"default:CURRENT_TIMESTAMP"`
	ApprovedBy   string    `json:"approved_by" gorm:"size:255;comment:Nama Pejabat yang menyetujui"`
	Reason       string    `json:"reason" gorm:"type:text;comment:Alasan mutasi"`

	// Relasi Struct
	Asset          Asset       `json:"asset" gorm:"foreignKey:AssetID"`
	FromDepartment *Department `json:"from_department" gorm:"foreignKey:FromDepartmentID"`
	ToDepartment   Department  `json:"to_department" gorm:"foreignKey:ToDepartmentID"`
	FromRoom       *Room       `json:"from_room" gorm:"foreignKey:FromRoomID"`
	ToRoom         Room        `json:"to_room" gorm:"foreignKey:ToRoomID"`
}
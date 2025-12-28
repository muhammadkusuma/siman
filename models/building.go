package models

import (
	"time"
)

// Building menyimpan data Gedung fisik.
type Building struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Code        string `json:"code" gorm:"size:50;comment:Kode Gedung (misal: G1)"`
	Name        string `json:"name" gorm:"size:255;comment:Nama Gedung (misal: Gedung Rektorat)"`
	TotalFloors int    `json:"total_floors" gorm:"comment:Jumlah Lantai dalam gedung"`

	// Relasi: Satu Gedung punya banyak Ruangan
	Rooms []Room `json:"rooms" gorm:"foreignKey:BuildingID"`
}

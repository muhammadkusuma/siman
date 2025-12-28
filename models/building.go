package models

// Building menyimpan data Gedung fisik.
type Building struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement;comment:ID Unik Gedung"`
	Code        string `json:"code" gorm:"size:50;comment:Kode Gedung Internal (cth: G1, GR)"`
	Name        string `json:"name" gorm:"size:255;comment:Nama Gedung (cth: Gedung Rektorat, Gedung Teater)"`
	TotalFloors int    `json:"total_floors" gorm:"comment:Jumlah Lantai dalam gedung (cth: 5)"`

	// Relasi
	Rooms []Room `json:"rooms" gorm:"foreignKey:BuildingID"`
}
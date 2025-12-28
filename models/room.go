package models

// Room menyimpan detail ruangan di dalam gedung.
type Room struct {
	ID         uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	BuildingID uint   `json:"building_id" gorm:"not null;comment:Foreign Key ke tabel Buildings"`
	
	RoomNumber string `json:"room_number" gorm:"size:50;comment:Nomor Pintu (misal: A.204)"`
	Name       string `json:"name" gorm:"size:255;comment:Fungsi Ruangan (misal: Lab Komputer)"`
	Floor      int    `json:"floor" gorm:"comment:Posisi Lantai"`

	// Relasi Belongs To
	Building Building `json:"building" gorm:"foreignKey:BuildingID"`
}
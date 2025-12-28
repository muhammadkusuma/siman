package models

// Room menyimpan detail ruangan di dalam gedung.
type Room struct {
	ID         uint   `json:"id" gorm:"primaryKey;autoIncrement;comment:ID Unik Ruangan"`
	BuildingID uint   `json:"building_id" gorm:"not null;comment:Foreign Key ke tabel Buildings"`
	
	RoomNumber string `json:"room_number" gorm:"size:50;comment:Nomor Pintu/Label Ruangan (cth: A.204, 101)"`
	Name       string `json:"name" gorm:"size:255;comment:Fungsi/Nama Ruangan (cth: Lab Komputer, Ruang Dosen)"`
	Floor      int    `json:"floor" gorm:"comment:Posisi Lantai (cth: 2)"`

	// Relasi
	Building Building `json:"building" gorm:"foreignKey:BuildingID"`
}
package models

import (
	"time"
)

// AssetCategory menyimpan kode klasifikasi BMN (Barang Milik Negara).
type AssetCategory struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	KodeBarang  string `json:"kode_barang" gorm:"size:50;unique;comment:Kode Klasifikasi BMN (misal: 3.10.01.02.001)"`
	Name        string `json:"name" gorm:"size:255;not null;comment:Nama Golongan Barang"`
	Description string `json:"description" gorm:"type:text;comment:Keterangan detail kategori"`
}
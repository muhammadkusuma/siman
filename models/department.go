package models

// Department merepresentasikan Prodi (S1/S2/S3) atau Bagian/Biro spesifik.
type Department struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement;comment:ID Unik Departemen"`
	FacultyID uint   `json:"faculty_id" gorm:"not null;comment:Foreign Key ke tabel Faculties"`
	
	Code      string `json:"code" gorm:"size:50;comment:Kode Internal Prodi/Unit (cth: IF-S1, TIF-D3)"`
	Name      string `json:"name" gorm:"size:255;not null;comment:Nama Prodi (cth: S1 Teknik Informatika)"`
	
	StudyLevel string `json:"study_level" gorm:"type:enum('D3','S1','S2','S3','Profesi','Non-Akademik');default:'Non-Akademik';comment:Jenjang Studi (cth: S1). Pilih Non-Akademik untuk kantor admin/biro."`
	
	HeadOfDepartment string `json:"head_of_department" gorm:"size:255;comment:Nama Kaprodi atau Kepala Unit saat ini (cth: Dr. Budi Santoso)"`

	// Relasi
	Faculty Faculty `json:"faculty" gorm:"foreignKey:FacultyID"`
}
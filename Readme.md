# SIMAN - Sistem Informasi Manajemen Aset

SIMAN adalah aplikasi Backend REST API yang dibangun menggunakan **Go (Golang)** dengan framework **Gin** dan ORM **GORM**. Aplikasi ini dirancang untuk mengelola inventaris aset organisasi (Kampus/Perusahaan), mencakup manajemen data master, pelacakan aset, riwayat perpindahan (mutasi), pencatatan perbaikan (maintenance), serta audit log aktivitas pengguna.

## 🚀 Fitur Utama

* **Otentikasi & Keamanan**: Login dan Register menggunakan **JWT (JSON Web Token)** dan hashing password dengan **Bcrypt**.
* **Manajemen Aset Lengkap**:
* CRUD Aset (Create, Read, Update, Delete).
* **Upload Foto Aset** (Multipart Form Data).
* Tracking otomatis user yang menginput (`CreatedBy`) dan mengupdate (`UpdatedBy`).


* **Manajemen Data Master**: Fakultas, Prodi (Departemen), Gedung, Ruangan, dan Kategori Aset.
* **Transaksi Aset**:
* **Mutasi**: Perpindahan aset antar unit/ruangan dengan riwayat lengkap.
* **Maintenance**: Pencatatan riwayat perbaikan aset, biaya, dan vendor.


* **Audit Logging**: Sistem secara otomatis mencatat setiap aktivitas sensitif (Siapa melakukan apa, IP address, perubahan data sebelum & sesudah) ke dalam database.
* **Static File Serving**: Akses langsung ke foto aset yang diupload.

## 🛠️ Tech Stack

* **Language**: Go (Golang) v1.20+
* **Framework**: Gin Web Framework
* **ORM**: GORM
* **Database**: MySQL
* **Authentication**: JWT-go
* **Utilities**: Crypto (Bcrypt)

## 📂 Struktur Project

```text
siman/
├── controllers/      # Logika bisnis (Handler API)
├── middlewares/      # Middleware (Auth JWT)
├── models/           # Definisi Struct & Konfigurasi DB
├── uploads/          # Folder penyimpanan foto aset (Auto-generated)
├── main.go           # Entry point & Routing
├── go.mod            # Dependency manager
└── README.md
```

## ⚙️ Instalasi & Konfigurasi

1. **Clone Repository**
```bash
git clone https://github.com/muhammadkusuma/siman.git
cd siman
```


2. **Install Dependencies**
```bash
go mod tidy
```


3. **Konfigurasi Database**
* Buat database MySQL bernama `db_siman`.
* Edit `models/setup.go` dan sesuaikan user/password:
```go
dsn := "root:12345678@tcp(127.0.0.1:3306)/db_siman?charset=utf8mb4&parseTime=True&loc=Local"
```




4. **Jalankan Aplikasi**
```bash
go run main.go
```


Server berjalan di `http://localhost:3000`.

---

## 📚 Dokumentasi API Lengkap

Semua endpoint di bawah grup `/api` membutuhkan Header:
`Authorization: Bearer <token_jwt>`

### 1. Otentikasi (Public)

| Method | Endpoint | Deskripsi | Body (JSON) |
| --- | --- | --- | --- |
| `POST` | `/register` | Pendaftaran user baru | `{"username": "admin", "password": "123", "full_name": "Admin", "role": "SuperAdmin", "email": "admin@test.com"}` |
| `POST` | `/login` | Login user & dapatkan Token | `{"username": "admin", "password": "123"}` |

### 2. Manajemen Aset (Assets)

Gunakan **Multipart/Form-Data** untuk Create & Update agar bisa upload foto.

| Method | Endpoint | Deskripsi | Parameter Form-Data |
| --- | --- | --- | --- |
| `GET` | `/api/assets` | Lihat semua aset | - |
| `GET` | `/api/assets/:id` | Detail satu aset | - |
| `POST` | `/api/assets` | Tambah aset + Foto | `name`, `inventory_code`, `nup`, `asset_category_id`, `department_id`, `room_id`, `price`, `photo` (File) |
| `PUT` | `/api/assets/:id` | Update data aset + Foto | (Sama seperti POST) |
| `DELETE` | `/api/assets/:id` | Hapus aset | - |

**Contoh CURL Create Asset:**

```bash
curl -X POST http://localhost:3000/api/assets \
  -H "Authorization: Bearer <TOKEN>" \
  -F "name=Laptop Dell" \
  -F "inventory_code=UIN-001" \
  -F "asset_category_id=1" \
  -F "department_id=1" \
  -F "room_id=1" \
  -F "photo=@C:\Gambar\laptop.jpg"
```

### 3. Transaksi & Log

| Method | Endpoint | Deskripsi | Body (JSON) |
| --- | --- | --- | --- |
| `GET` | `/api/mutations` | Riwayat perpindahan aset | - |
| `POST` | `/api/mutations` | Pindahkan aset (Mutasi) | `{"asset_id": 1, "to_department_id": 2, "to_room_id": 5, "reason": "Pindah Gedung", "approved_by": "Kabag"}` |
| `GET` | `/api/maintenances` | Riwayat perbaikan | - |
| `POST` | `/api/maintenances` | Lapor kerusakan | `{"asset_id": 1, "issue_date": "2024-01-01T00:00:00Z", "description": "Layar Mati", "status": "Pending"}` |
| `GET` | `/api/audit-logs` | Lihat Log Aktivitas User | - |

### 4. Data Master (Referensi)

Endpoint ini digunakan untuk mengisi dropdown atau data referensi.

| Method | Endpoint | Deskripsi | Body (JSON) |
| --- | --- | --- | --- |
| `GET` | `/api/faculties` | Daftar Fakultas | - |
| `POST` | `/api/faculties` | Tambah Fakultas | `{"code": "FST", "name": "Sains Tek", "type": "Fakultas"}` |
| `GET` | `/api/departments` | Daftar Prodi/Unit | - |
| `POST` | `/api/departments` | Tambah Prodi | `{"faculty_id": 1, "code": "IF", "name": "Informatika", "study_level": "S1"}` |
| `GET` | `/api/buildings` | Daftar Gedung | - |
| `POST` | `/api/buildings` | Tambah Gedung | `{"code": "G1", "name": "Gedung Kuliah", "total_floors": 3}` |
| `POST` | `/api/rooms` | Tambah Ruangan | `{"building_id": 1, "room_number": "101", "name": "Lab", "floor": 1}` |
| `GET` | `/api/categories` | Kategori Aset (BMN) | - |
| `POST` | `/api/categories` | Tambah Kategori | `{"kode_barang": "3.10.01", "name": "Laptop"}` |

### 5. Akses File Statis

Foto yang diupload dapat diakses langsung melalui browser:
`http://localhost:3000/uploads/assets/<nama_file.jpg>`
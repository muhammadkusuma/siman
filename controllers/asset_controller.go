package controllers

import (
	"fmt"
	"muhammadkusuma/siman/models"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// Update fungsi GetAssets agar bisa search
func GetAssets(c *gin.Context) {
	var assets []models.Asset

	// Ambil parameter query dari URL: ?search=laptop&status=Baik
	search := c.Query("search")
	status := c.Query("status")

	// Mulai Query
	query := models.DB.Preload("Category").Preload("Department").Preload("Room")

	// Filter by Name (Search)
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Filter by Condition
	if status != "" {
		query = query.Where("condition_status = ?", status)
	}

	// Eksekusi
	query.Find(&assets)

	c.JSON(http.StatusOK, gin.H{"data": assets})
}

// GetAssetByID (Tidak berubah)
func GetAssetByID(c *gin.Context) {
	var asset models.Asset
	if err := models.DB.Preload("Category").Preload("Department").Preload("Room").Preload("CreatedBy").Preload("UpdatedBy").First(&asset, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": asset})
}

// Struct Input untuk form-data agar parsing tanggal aman
type AssetForm struct {
	Name              string  `form:"name"`
	AssetCategoryID   uint    `form:"asset_category_id"`
	InventoryCode     string  `form:"inventory_code"`
	NUP               int     `form:"nup"`
	Brand             string  `form:"brand"`
	Model             string  `form:"model"`
	SerialNumber      string  `form:"serial_number"`
	DepartmentID      uint    `form:"department_id"`
	RoomID            uint    `form:"room_id"`
	AcquisitionDate   string  `form:"acquisition_date"` // Terima string, parsing manual
	Price             float64 `form:"price"`
	ConditionStatus   string  `form:"condition_status"`
	OperationalStatus string  `form:"operational_status"`
}

func CreateAsset(c *gin.Context) {
	var form AssetForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Manual Mapping ke Model
	input := models.Asset{
		Name:              form.Name,
		AssetCategoryID:   form.AssetCategoryID,
		InventoryCode:     form.InventoryCode,
		NUP:               form.NUP,
		Brand:             form.Brand,
		Model:             form.Model,
		SerialNumber:      form.SerialNumber,
		DepartmentID:      form.DepartmentID,
		RoomID:            form.RoomID,
		Price:             form.Price,
		ConditionStatus:   form.ConditionStatus,
		OperationalStatus: form.OperationalStatus,
	}

	// Parsing Tanggal Manual
	if form.AcquisitionDate != "" {
		t, err := time.Parse(time.RFC3339, form.AcquisitionDate) // Parse ISO format dari JS
		if err == nil {
			input.AcquisitionDate = t
		}
	}

	// Upload File Logic
	file, err := c.FormFile("photo")
	if err == nil {
		path := "uploads/assets"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.MkdirAll(path, os.ModePerm)
		}
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
		filepath := fmt.Sprintf("%s/%s", path, filename)
		if err := c.SaveUploadedFile(file, filepath); err == nil {
			input.PhotoPath = filepath
		}
	}

	// Set User
	if userID, exists := c.Get("userID"); exists {
		input.CreatedByID = userID.(uint)
		input.UpdatedByID = userID.(uint)
	}

	// Simpan
	if err := models.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan: " + err.Error()})
		return
	}

	RecordLog(c, "CREATE", "assets", input.ID, "Menambahkan aset baru: "+input.Name)
	c.JSON(http.StatusOK, gin.H{"data": input})
}

func UpdateAsset(c *gin.Context) {
	var asset models.Asset
	if err := models.DB.First(&asset, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found!"})
		return
	}

	var form AssetForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update Fields
	asset.Name = form.Name
	asset.AssetCategoryID = form.AssetCategoryID
	asset.InventoryCode = form.InventoryCode
	asset.NUP = form.NUP
	asset.Brand = form.Brand
	asset.Model = form.Model
	asset.SerialNumber = form.SerialNumber
	asset.DepartmentID = form.DepartmentID
	asset.RoomID = form.RoomID
	asset.Price = form.Price
	asset.ConditionStatus = form.ConditionStatus
	asset.OperationalStatus = form.OperationalStatus

	if form.AcquisitionDate != "" {
		t, err := time.Parse(time.RFC3339, form.AcquisitionDate)
		if err == nil {
			asset.AcquisitionDate = t
		}
	}

	// Upload File Logic (Update)
	file, err := c.FormFile("photo")
	if err == nil {
		path := "uploads/assets"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.MkdirAll(path, os.ModePerm)
		}
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
		newFilePath := fmt.Sprintf("%s/%s", path, filename)

		if err := c.SaveUploadedFile(file, newFilePath); err == nil {
			if asset.PhotoPath != "" {
				os.Remove(asset.PhotoPath) // Hapus foto lama
			}
			asset.PhotoPath = newFilePath
		}
	}

	if userID, exists := c.Get("userID"); exists {
		asset.UpdatedByID = userID.(uint)
	}

	// Simpan & Cek Error
	if err := models.DB.Save(&asset).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update: " + err.Error()})
		return
	}

	RecordLog(c, "UPDATE", "assets", asset.ID, "Mengupdate data aset")
	c.JSON(http.StatusOK, gin.H{"data": asset})
}

// DeleteAsset (Tidak berubah)
func DeleteAsset(c *gin.Context) {
	var asset models.Asset
	if err := models.DB.First(&asset, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found!"})
		return
	}

	// Hapus file fisik foto jika ada
	if asset.PhotoPath != "" {
		os.Remove(asset.PhotoPath)
	}

	models.DB.Delete(&asset)
	RecordLog(c, "DELETE", "assets", asset.ID, "Menghapus aset: "+asset.Name)
	c.JSON(http.StatusOK, gin.H{"data": true})
}

package controllers

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"news/models"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ShowNews(c *fiber.Ctx) error {
	var data []models.News

	if err := models.DB.Table("news").
		Select("news.*, category.title AS category_name").
		Joins("left join category on category.id = news.category").
		Find(&data).Error; err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	if len(data) == 0 {
		return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
	}
	return c.JSON(data)
}

/* func ShowNews(c *fiber.Ctx) error {
    var data []struct {
        Image        string `json:"image"`
        CategoryName string `json:"category_name"`
    }

    // Perform the join between news and category tables
    if err := models.DB.Table("news").
        Select("news.image, category.title as category_name").
        Joins("left join category on news.category = category.id").
        Find(&data).Error; err != nil {
        return c.JSON(fiber.Map{
            "status":  fiber.StatusInternalServerError,
            "message": "Failed to load data",
            "error":   err.Error(),
        })
    }

    if len(data) == 0 {
        return c.JSON(fiber.Map{
            "status":  fiber.StatusNotFound,
            "message": "No data found",
        })
    }

    return c.JSON(fiber.Map{
        "status":  fiber.StatusOK,
        "message": "Data loaded successfully",
        "data":    data,
    })
} */

func IndexNews(c *fiber.Ctx) error {
	id := c.Params("id")
	var data models.News

	if err := models.DB.First(&data, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
		}
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	return c.JSON(data)
}

func CreateNews(c *fiber.Ctx) error {
	var data models.News

	if err := c.BodyParser(&data); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	if err := models.DB.Create(&data).Error; err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to save data", err.Error())
	}

	return jsonResponse(c, fiber.StatusCreated, "Data successfully added", data)
}

func UpdateNews(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil)
	}

	var data models.News
	if err := models.DB.First(&data, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
		}
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	var updateData models.News
	if err := c.BodyParser(&updateData); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	if updateData.ID != 0 && updateData.ID != id {
		if err := models.DB.First(&models.News{}, updateData.ID).Error; err == nil {
			return jsonResponse(c, fiber.StatusBadRequest, "The updated ID is already in use", nil)
		}
	}

	if err := models.DB.Model(&data).Updates(updateData).Error; err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to update data", err.Error())
	}

	return jsonResponse(c, fiber.StatusOK, "Data successfully updated", nil)
}

func DeleteNews(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validasi ID
	if id == "" {
		return jsonResponse(c, fiber.StatusBadRequest, "id is required", nil)
	}

	// Cek apakah berita dengan ID tersebut ada
	var news models.News
	if err := models.DB.First(&news, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return jsonResponse(c, fiber.StatusNotFound, "data not found or already deleted", nil)
		}
		fmt.Println("database error:", err)
		return jsonResponse(c, fiber.StatusInternalServerError, "database error", nil)
	}

	// Direktori upload
	uploadDir := "/var/www/html/images/garuda/news"
	filePath := filepath.Join(uploadDir, fmt.Sprintf("%s.jpg", id))

	// Hapus file jika ada
	if _, err := os.Stat(filePath); err == nil {
		// File exists, attempt to delete
		if err := os.Remove(filePath); err != nil {
			fmt.Printf("error deleting file (%s): %v\n", filePath, err)
			return jsonResponse(c, fiber.StatusInternalServerError, "failed to delete associated file", nil)
		}
		fmt.Printf("file deleted: %s\n", filePath)
	} else if !os.IsNotExist(err) {
		// Error selain file tidak ada
		fmt.Printf("error checking file (%s): %v\n", filePath, err)
		return jsonResponse(c, fiber.StatusInternalServerError, "failed to check associated file", nil)
	}

	// Hapus data berita dari database
	if models.DB.Delete(&models.News{}, id).RowsAffected == 0 {
		return jsonResponse(c, fiber.StatusNotFound, "data not found or already deleted", nil)
	}

	return jsonResponse(c, fiber.StatusOK, "successfully deleted data and associated file", nil)
}

// Fungsi utama upload image
func UploadNewsImage(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validasi ID
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id is required",
		})
	}

	// Parsing ID ke tipe integer
	newsID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}

	// Ambil file dari request
	file, err := c.FormFile("image")
	if err != nil {
		fmt.Println("error reading file:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to read the file. ensure 'image' key is included in the form-data request",
		})
	}

	fmt.Printf("File received: %s (Size: %d bytes)\n", file.Filename, file.Size)

	// Validasi ukuran file (maks 5 MB)
	const maxFileSize = 5 * 1024 * 1024 // 5 MB
	if file.Size > maxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file size exceeds 5 MB limit",
		})
	}

	// Validasi tipe file dengan MIME
	if err := validateFileType(file); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Cek apakah berita dengan ID tersebut ada
	var news models.News
	if err := models.DB.First(&news, "id = ?", newsID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "news with the given id not found",
			})
		}
		fmt.Println("database error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "database error",
		})
	}

	// Direktori upload
	uploadDir := "/var/www/html/images/garuda/news"
	if err := ensureDirectoryExists(uploadDir); err != nil {
		fmt.Println("error ensuring directory exists:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create upload directory",
		})
	}

	// Sanitasi nama file
	ext := filepath.Ext(file.Filename)
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	if !allowedExtensions[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unsupported file type",
		})
	}

	fileName := fmt.Sprintf("%d%s", newsID, ext)
	filePath := filepath.Join(uploadDir, fileName)

	// Pastikan path tetap di dalam direktori yang diizinkan
	if !strings.HasPrefix(filepath.Clean(filePath), filepath.Clean(uploadDir)) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "invalid file path",
		})
	}

	// Simpan file
	if err := c.SaveFile(file, filePath); err != nil {
		fmt.Printf("error saving file (%s): %v\n", filePath, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save the file",
		})
	}

	// URL publik untuk file
	publicURL := fmt.Sprintf("https://web.ayomenjadi.com/images/garuda/news/%s", fileName)

	// Perbarui data di database
	news.Image = publicURL
	news.UpdatedAt = time.Now()
	if err := models.DB.Save(&news).Error; err != nil {
		fmt.Println("error updating database:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update news image",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "image uploaded successfully",
		"image_path": publicURL,
	})
}

// Fungsi untuk memastikan direktori ada
func ensureDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// Fungsi untuk validasi tipe file menggunakan MIME
func validateFileType(file *multipart.FileHeader) error {
	fileHeader, err := file.Open()
	if err != nil {
		return fmt.Errorf("unable to open file")
	}
	defer fileHeader.Close()

	buffer := make([]byte, 512)
	if _, err := fileHeader.Read(buffer); err != nil {
		return fmt.Errorf("unable to read file header")
	}

	// Deteksi tipe MIME
	mimeType := http.DetectContentType(buffer)
	allowedMimeTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}
	if !allowedMimeTypes[mimeType] {
		return fmt.Errorf("invalid file type")
	}
	return nil
}


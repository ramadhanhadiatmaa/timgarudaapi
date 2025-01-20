package controllers

import (
	"errors"
	"fmt"
	"news/models"
	"os"
	"path/filepath"
	"strconv"
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

	if models.DB.Delete(&models.News{}, id).RowsAffected == 0 {
		return jsonResponse(c, fiber.StatusNotFound, "Data not found or already deleted", nil)
	}

	return jsonResponse(c, fiber.StatusOK, "Successfully deleted data", nil)
}

func UploadNewsImage(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validasi ID
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	// Parsing ID ke tipe integer
	newsID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	// Ambil file dari request
	file, err := c.FormFile("image_path")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to read the file",
		})
	}

	fmt.Printf("File received: %s (Size: %d bytes)\n", file.Filename, file.Size)

	// Cek apakah berita dengan ID tersebut ada
	var news models.News
	if err := models.DB.First(&news, "id = ?", newsID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "News with the given ID not found",
			})
		}
		fmt.Println("Database error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Direktori upload
	uploadDir := "/var/www/html/images/garuda/news"
	if err := ensureDirectoryExists(uploadDir); err != nil {
		fmt.Println("Error ensuring directory exists:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create upload directory",
		})
	}

	// Buat nama file berdasarkan ID
	ext := filepath.Ext(file.Filename)
	safeFileName := sanitizeFilename(file.Filename)
	fileName := fmt.Sprintf("%d%s", newsID, ext)
	filePath := filepath.Join(uploadDir, safeFileName)

	// Simpan file
	if err := c.SaveFile(file, filePath); err != nil {
		fmt.Println("Error saving file:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save the file",
		})
	}

	// URL publik untuk file
	publicURL := fmt.Sprintf("https://web.ayomenjadi.com/images/garuda/news/%s", fileName)

	// Perbarui data di database
	news.Image = publicURL
	news.UpdatedAt = time.Now()
	if err := models.DB.Save(&news).Error; err != nil {
		fmt.Println("Error updating database:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update news image",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Image uploaded successfully",
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

// Fungsi untuk sanitasi nama file
func sanitizeFilename(filename string) string {
	return filepath.Base(filename)
}

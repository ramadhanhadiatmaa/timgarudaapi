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

	file, err := c.FormFile("image_path")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to read the file",
		})
	}

	fmt.Println("File received:", file.Filename, "Size:", file.Size)

	var news models.News
	if err := models.DB.First(&news, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "News with the given ID not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	uploadDir := "/var/www/html/images/garuda/news"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			fmt.Println("Error creating directory:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Unable to create upload directory",
			})
		}
	}

	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s%s", id, ext)
	filePath := filepath.Join(uploadDir, fileName)

	if err := c.SaveFile(file, filePath); err != nil {
		fmt.Println("Error saving file:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save the file",
		})
	}

	publicURL := fmt.Sprintf("https://web.ayomenjadi.com/images/garuda/news/%s", fileName)

	news.Image = publicURL
	news.UpdatedAt = time.Now()
	if err := models.DB.Save(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update news image",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Image uploaded successfully",
		"image_path": publicURL,
	})
}

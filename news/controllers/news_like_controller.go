package controllers

import (
	"news/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ShowNewsLike(c *fiber.Ctx) error {
	var data []models.NewsLike

	if err := models.DB.Find(&data).Error; err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	if len(data) == 0 {
		return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
	}

	return c.JSON(data)
}

func IndexNewsLike(c *fiber.Ctx) error {
	id := c.Params("id")
	var data models.NewsLike

	if err := models.DB.First(&data, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
		}
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	return c.JSON(data)
}

func CreateNewsLike(c *fiber.Ctx) error {
	var data models.NewsLike

	if err := c.BodyParser(&data); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	if err := models.DB.Create(&data).Error; err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to save data", err.Error())
	}

	return jsonResponse(c, fiber.StatusCreated, "Data successfully added", data)
}

func DeleteNewsLike(c *fiber.Ctx) error {
	id := c.Params("id")

	if models.DB.Delete(&models.NewsLike{}, id).RowsAffected == 0 {
		return jsonResponse(c, fiber.StatusNotFound, "Data not found or already deleted", nil)
	}

	return jsonResponse(c, fiber.StatusOK, "Successfully deleted data", nil)
}
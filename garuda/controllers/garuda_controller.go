package controllers

import (
	"garuda/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Index(c *fiber.Ctx) error {
	id := c.Params("id")
	var data models.Garuda

	if err := models.DB.First(&data, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
		}
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	return c.JSON(data)
}

func Create(c *fiber.Ctx) error {
	var data models.Garuda

	if err := c.BodyParser(&data); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	if err := models.DB.Create(&data).Error; err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to save data", err.Error())
	}

	return jsonResponse(c, fiber.StatusCreated, "Data successfully added", data)
}

func Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil)
	}

	var data models.Garuda
	if err := models.DB.First(&data, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
		}
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	var updateData models.Garuda
	if err := c.BodyParser(&updateData); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	if updateData.ID != 0 && updateData.ID != id {
		if err := models.DB.First(&models.Garuda{}, updateData.ID).Error; err == nil {
			return jsonResponse(c, fiber.StatusBadRequest, "The updated ID is already in use", nil)
		}
	}

	if err := models.DB.Model(&data).Updates(updateData).Error; err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to update data", err.Error())
	}

	return jsonResponse(c, fiber.StatusOK, "Data successfully updated", nil)
}

func Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if models.DB.Delete(&models.Garuda{}, id).RowsAffected == 0 {
		return jsonResponse(c, fiber.StatusNotFound, "Data not found or already deleted", nil)
	}

	return jsonResponse(c, fiber.StatusOK, "Successfully deleted data", nil)
}

func jsonResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"data":    data,
	})
}
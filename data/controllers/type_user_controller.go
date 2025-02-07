package controllers

import (
	"data/models"

	"github.com/gofiber/fiber/v2"
)

func CreateType(c *fiber.Ctx) error {
	var data models.TypeUser

	if err := c.BodyParser(&data); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	if err := models.DB.Create(&data).Error; err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to save data", err.Error())
	}

	return jsonResponse(c, fiber.StatusCreated, "Data successfully added", data)
}

func DeleteType(c *fiber.Ctx) error {
	id := c.Params("id")

	if models.DB.Delete(&models.TypeUser{}, id).RowsAffected == 0 {
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
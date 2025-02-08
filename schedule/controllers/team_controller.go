package controllers

import (
	"fmt"
	"schedule/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ShowTeam(c *fiber.Ctx) error {
	var data []models.Team

	if err := models.DB.Find(&data).Error; err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	if len(data) == 0 {
		return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
	}

	return c.JSON(data)
}

func IndexTeam(c *fiber.Ctx) error {
	id := c.Params("id")
	var data models.Team

	if err := models.DB.First(&data, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
		}
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	return c.JSON(data)
}

func CreateTeam(c *fiber.Ctx) error {
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	allowedKeys := []string{"team", "flag"}

	for key := range data {
		if !contains(allowedKeys, key) {
			return jsonResponse(c, fiber.StatusBadRequest, "Inputting data is not allowed", nil)
		}
	}

	if exampleValue, exists := data["team"]; exists {
		dataVal := models.Team{
			Team: fmt.Sprintf("%v", exampleValue),
		}

		if err := models.DB.Create(&dataVal).Error; err != nil {
			return jsonResponse(c, fiber.StatusInternalServerError, "Failed to save data", err.Error())
		}

		return jsonResponse(c, fiber.StatusCreated, "Data successfully added", dataVal)
	}
	
	return jsonResponse(c, fiber.StatusBadRequest, "'example' key is required", nil)
}

func UpdateTeam(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil)
	}

	var data models.Team
	if err := models.DB.First(&data, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
		}
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	var updateData models.Team
	if err := c.BodyParser(&updateData); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	if updateData.ID != 0 && updateData.ID != id {
		if err := models.DB.First(&models.Team{}, updateData.ID).Error; err == nil {
			return jsonResponse(c, fiber.StatusBadRequest, "The updated ID is already in use", nil)
		}
	}

	if err := models.DB.Model(&data).Updates(updateData).Error; err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to update data", err.Error())
	}

	return jsonResponse(c, fiber.StatusOK, "Data successfully updated", nil)
}

func DeleteTeam(c *fiber.Ctx) error {
	id := c.Params("id")

	if models.DB.Delete(&models.Team{}, id).RowsAffected == 0 {
		return jsonResponse(c, fiber.StatusNotFound, "Data not found or already deleted", nil)
	}

	return jsonResponse(c, fiber.StatusOK, "Successfully deleted data", nil)
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

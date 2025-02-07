package controllers

import (
	"schedule/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Show(c *fiber.Ctx) error {
	var data []models.Schedule
	if err := models.DB.Preload("Home").Preload("Away").
		Find(&data).Error; err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load schedules", err.Error())
	}

	if len(data) == 0 {
		return jsonResponse(c, fiber.StatusNotFound, "No schedules found", nil)
	}

	var schedule models.Schedule

	return c.JSON(fiber.Map{
		"id":        schedule.ID,
		"home":      schedule.Home.Team,
		"away":      schedule.Away.Team,
		"flag_home": schedule.Home.Flag,
		"flag_away": schedule.Away.Flag,
		"venue":     schedule.Venue,
		"date":      schedule.Date,
	})
}

func Index(c *fiber.Ctx) error {
	id := c.Params("id")
	var data models.Schedule

	// Load data dengan preload Home dan Away
	if err := models.DB.Preload("Home").Preload("Away").
		First(&data, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
		}
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	var schedule models.Schedule

	return c.JSON(fiber.Map{
		"id":        schedule.ID,
		"home":      schedule.Home.Team,
		"away":      schedule.Away.Team,
		"flag_home": schedule.Home.Flag,
		"flag_away": schedule.Away.Flag,
		"venue":     schedule.Venue,
		"date":      schedule.Date,
	})
}

func Create(c *fiber.Ctx) error {
	var data models.Schedule

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

	var data models.Schedule
	if err := models.DB.First(&data, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
		}
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	var updateData models.Schedule
	if err := c.BodyParser(&updateData); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	if updateData.ID != 0 && updateData.ID != id {
		if err := models.DB.First(&models.Schedule{}, updateData.ID).Error; err == nil {
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

	if models.DB.Delete(&models.Schedule{}, id).RowsAffected == 0 {
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

package controllers

import (
	"auth/models"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	typeUser, err := strconv.Atoi(data["type"])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Type User"})
	}

	var existingUser models.User
	if err := models.DB.First(&existingUser, "email = ?", data["email"]).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already exists"})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		Email:     data["email"],
		Password:  string(password),
		FullName:  data["full_name"],
		Phone:     data["phone"],
		Type:      typeUser,
		CreatedAt: time.Now(),
	}

	if err := models.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not register user"})
	}

	return c.JSON(fiber.Map{"message": "User registered successfully"})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	// Parse input data
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	// Use Joins to explicitly perform a LEFT JOIN between `user` and `type_user`
	err := models.DB.Debug().Joins("LEFT JOIN type_user ON type_user.id = user.type").
		Where("email = ?", data["email"]).
		First(&user).Error

	// If there's an error, return unauthorized
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// Print user details to verify if the data is correct
	fmt.Printf("User: %+v\n", user)

	// Compare password hash with input password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])) != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// Generate JWT token
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Secret key not configured"})
	}

	// Create JWT claims
	claims := jwt.MapClaims{
		"email": user.Email,
		"type":  user.Type,
		"exp":   time.Now().Add(time.Hour * 240).Unix(),
	}

	// Create the token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token and get the result
	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	// Print TypeInfo to verify if it was preloaded correctly
	fmt.Printf("TypeInfo: %+v\n", user.TypeInfo)

	// Return the response with the token and user info
	return c.JSON(fiber.Map{
		"token":     t,
		"email":     user.Email,
		"phone":     user.Phone,
		"full_name": user.FullName,
		"type":      user.TypeInfo.Type, // This should now show the correct type from TypeUser table
	})
}

func Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil)
	}

	var data models.User
	if err := models.DB.First(&data, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
		}
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	var updateData models.User
	if err := c.BodyParser(&updateData); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid", err.Error())
	}

	if updateData.ID != 0 && updateData.ID != id {
		if err := models.DB.First(&models.User{}, updateData.ID).Error; err == nil {
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

	if models.DB.Delete(&models.User{}, id).RowsAffected == 0 {
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

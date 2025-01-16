package controllers

import (
	"auth/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	if err := models.DB.Preload("TypeInfo").First(&user, "email = ?", data["email"]).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])) != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Secret key not configured"})
	}

	claims := jwt.MapClaims{
		"email": user.Email,
		"type":  user.Type,
		"exp":   time.Now().Add(time.Hour * 240).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.JSON(fiber.Map{
		"token":     t,
		"email":     user.Email,
		"phone":     user.Phone,
		"full_name": user.FullName,
		"type":      user.TypeInfo.Type,
	})
}

func Update(c *fiber.Ctx) error {
	email := c.Params("email")
	var user models.User

	if err := models.DB.First(&user, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return jsonResponse(c, fiber.StatusNotFound, "User not found", nil)
		}
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load user", err.Error())
	}

	if err := c.BodyParser(&user); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid input data", nil)
	}

	if err := models.DB.Save(&user).Error; err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to update user", err.Error())
	}

	return jsonResponse(c, fiber.StatusOK, "User updated successfully", user)
}

func UploadUserImage(c *fiber.Ctx) error {
	email := c.Params("email")

	// Retrieve the file from the form
	file, err := c.FormFile("image_path")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to read the file",
		})
	}

	fmt.Println("File received:", file.Filename, "Size:", file.Size)

	// Validate user exists
	var user models.User
	if err := models.DB.First(&user, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Save the file to the specified directory
	uploadDir := "/var/www/html/images" // Local directory mapped to Docker container
	fmt.Println("Saving file to directory:", uploadDir)

	// Ensure the directory exists
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Unable to create directory",
			})
		}
	}

	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s%s", email, ext)
	filePath := filepath.Join(uploadDir, fileName)
	fmt.Println("Saving file to:", filePath)

	fileContent, err := file.Open()
	if err != nil {
		fmt.Println("Error opening file:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to open the file",
		})
	}
	defer fileContent.Close()

	outFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to create the file on the server",
		})
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, fileContent)
	if err != nil {
		fmt.Println("Error saving file content:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to save the file content",
		})
	}

	// Construct the public URL for the image
	publicURL := fmt.Sprintf("https://web.ayomenjadi.com/images/%s", fileName)

	// Update user information in the database
	user.Password = publicURL
	user.UpdatedAt = time.Now()
	if err := models.DB.Save(&user).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to update user information",
		})
	}

	// Return success response with image URL
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":    "Image uploaded successfully",
		"image_path": publicURL,
	})
}

func Delete(c *fiber.Ctx) error {
	email := c.Params("email")

	var user models.User
	if err := models.DB.First(&user, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return jsonResponse(c, fiber.StatusNotFound, "Data not found", nil)
		}
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	if err := models.DB.Delete(&user).Error; err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to delete data", err.Error())
	}

	return jsonResponse(c, fiber.StatusOK, "Successfully deleted data", nil)
}

func jsonResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"data":    data,
	})
}

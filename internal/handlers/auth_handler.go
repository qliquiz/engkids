package handlers

import (
	"engkids/internal/models"
	"engkids/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user by providing email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User Register"
// @Success 200 {object} models.User
// @Failure 400 {object} fiber.Map
// @Router /api/register [post]
func Register(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body models.User
		if err := c.BodyParser(&body); err != nil {
			return err
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
		user := models.User{Email: body.Email, Password: string(hash)}
		if err := db.Create(&user).Error; err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Email already in use"})
		}
		return c.JSON(user)
	}
}

// Login godoc
// @Summary Login a user
// @Description Login a user with email and password to get a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User Login"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Router /api/login [post]
func Login(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body models.User
		if err := c.BodyParser(&body); err != nil {
			return err
		}
		var user models.User
		db.First(&user, "email = ?", body.Email)
		if user.ID == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid credentials"})
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Wrong password"})
		}
		token, _ := jwt.GenerateToken(user.ID)
		return c.JSON(fiber.Map{"token": token})
	}
}

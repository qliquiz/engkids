package handlers

import (
	//"engkids/internal/models"
	"engkids/pkg/elasticsearch"
	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	//"github.com/sirupsen/logrus"
	//"golang.org/x/crypto/bcrypt"
	//"gorm.io/gorm"
)

/*

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
func Register(db *gorm.DB, log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body models.User
		if err := c.BodyParser(&body); err != nil {
			log.WithError(err).Error("Failed to parse registration request")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		log.WithField("email", body.Email).Info("Registration attempt")

		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
		if err != nil {
			log.WithError(err).Error("Failed to hash password")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}

		user := models.User{Email: body.Email, Password: string(hash)}
		if err := db.Create(&user).Error; err != nil {
			log.WithFields(logrus.Fields{
				"email": body.Email,
				"error": err.Error(),
			}).Error("User registration failed")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email already in use"})
		}

		log.WithFields(logrus.Fields{
			"user_id": user.ID,
			"email":   user.Email,
		}).Info("User registered successfully")

		// Не возвращаем пароль в ответе
		user.Password = ""
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
func Login(db *gorm.DB, log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body models.User
		if err := c.BodyParser(&body); err != nil {
			log.WithError(err).Error("Failed to parse login request")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		log.WithField("email", body.Email).Info("Login attempt")

		var user models.User
		db.First(&user, "email = ?", body.Email)
		if user.ID == 0 {
			log.WithField("email", body.Email).Warn("Login failed: user not found")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
			log.WithFields(logrus.Fields{
				"user_id": user.ID,
				"email":   user.Email,
			}).Warn("Login failed: wrong password")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Wrong password"})
		}

		token, err := jwt.GenerateToken(user.ID)
		if err != nil {
			log.WithError(err).Error("Failed to generate token")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
		}

		log.WithFields(logrus.Fields{
			"user_id": user.ID,
			"email":   user.Email,
		}).Info("User logged in successfully")

		return c.JSON(fiber.Map{"token": token})
	}
}
*/

// GetLogs godoc
// @Summary Get application logs
// @Description Get application logs from Elasticsearch
// @Tags logs
// @Accept json
// @Produce json
// @Param from query string false "Start time (ISO8601 format)"
// @Param to query string false "End time (ISO8601 format)"
// @Param level query string false "Log level (info, warn, error)"
// @Success 200 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/logs [get]
func GetLogs(es *elasticsearch.Client, log *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		from := c.Query("from")
		to := c.Query("to")
		level := c.Query("level")

		// Формируем запрос к Elasticsearch
		query := map[string]interface{}{
			"size": 100,
			"sort": []map[string]interface{}{
				{
					"@timestamp": map[string]string{
						"order": "desc",
					},
				},
			},
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{},
				},
			},
		}

		// Добавляем условия по времени, если указаны
		if from != "" || to != "" {
			timeQuery := map[string]interface{}{
				"range": map[string]interface{}{
					"@timestamp": map[string]interface{}{},
				},
			}

			if from != "" {
				timeQuery["range"].(map[string]interface{})["@timestamp"].(map[string]interface{})["gte"] = from
			}

			if to != "" {
				timeQuery["range"].(map[string]interface{})["@timestamp"].(map[string]interface{})["lte"] = to
			}

			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
				query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
				timeQuery,
			)
		}

		// Добавляем условие по уровню лога, если указан
		if level != "" {
			levelQuery := map[string]interface{}{
				"match": map[string]interface{}{
					"level": level,
				},
			}

			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
				query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
				levelQuery,
			)
		}

		// Отправляем запрос в Elasticsearch
		result, err := es.Search("engkids-logs-*", query)
		if err != nil {
			log.WithError(err).Error("Failed to search logs in Elasticsearch")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve logs",
			})
		}

		log.Info("Logs retrieved successfully")

		return c.JSON(result)
	}
}

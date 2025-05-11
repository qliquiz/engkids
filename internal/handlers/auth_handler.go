package handlers

import (
	"engkids/internal/dto"
	"engkids/internal/errors"
	"engkids/internal/services"
	"engkids/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	Service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := utils.ParseAndValidate(c, &req); err != nil {
		return errors.Handle(c, err)
	}

	resp, err := h.Service.Register(&req)
	if err != nil {
		return errors.Handle(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := utils.ParseAndValidate(c, &req); err != nil {
		return errors.Handle(c, err)
	}

	resp, err := h.Service.Login(&req)
	if err != nil {
		return errors.Handle(c, err)
	}

	return c.JSON(resp)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if err := utils.ParseAndValidate(c, &body); err != nil {
		return errors.Handle(c, err)
	}

	resp, err := h.Service.Refresh(body.RefreshToken)
	if err != nil {
		return errors.Handle(c, err)
	}

	return c.JSON(resp)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if err := utils.ParseAndValidate(c, &body); err != nil {
		return errors.Handle(c, err)
	}

	if err := h.Service.Logout(body.RefreshToken); err != nil {
		return errors.Handle(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

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
/*func GetLogs(es *elasticsearch.Client, log *logrus.Logger) fiber.Handler {
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
		result, err := es.Search("filebeat-*", query)
		if err != nil {
			log.WithError(err).Error("Failed to search logs in Elasticsearch")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve logs",
			})
		}

		log.Info("Logs retrieved successfully")

		return c.JSON(result)
	}
}*/

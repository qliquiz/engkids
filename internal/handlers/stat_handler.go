package handlers

import (
	"engkids/internal/dto"
	"engkids/internal/errors"
	"engkids/internal/services"
	"engkids/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// UserHandler обработчик запросов пользователя
type UserHandler struct {
	Service *services.UserService
}

// NewUserHandler создает новый обработчик
func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

// GetUserProfile получает профиль пользователя
func (h *UserHandler) GetUserProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	
	profile, err := h.Service.GetUserProfile(userID)
	if err != nil {
		return errors.Handle(c, err)
	}
	
	return c.JSON(profile)
}

// GetUserInventory получает инвентарь пользователя
func (h *UserHandler) GetUserInventory(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	
	inventory, err := h.Service.GetUserInventory(userID)
	if err != nil {
		return errors.Handle(c, err)
	}
	
	return c.JSON(fiber.Map{
		"inventory": inventory,
	})
}

// UpdateInventoryItem обновляет предмет в инвентаре
func (h *UserHandler) UpdateInventoryItem(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	
	var req dto.UpdateInventoryRequest
	if err := utils.ParseAndValidate(c, &req); err != nil {
		return errors.Handle(c, err)
	}
	
	if err := h.Service.UpdateInventoryItem(userID, &req); err != nil {
		return errors.Handle(c, err)
	}
	
	return c.JSON(fiber.Map{
		"status": "success",
		"message": "Предмет успешно обновлен",
	})
}

// PurchaseItem покупает предмет
func (h *UserHandler) PurchaseItem(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	
	var req dto.PurchaseItemRequest
	if err := utils.ParseAndValidate(c, &req); err != nil {
		return errors.Handle(c, err)
	}
	
	if err := h.Service.PurchaseItem(userID, &req); err != nil {
		return errors.Handle(c, err)
	}
	
	// Получаем обновленный профиль после покупки
	profile, err := h.Service.GetUserProfile(userID)
	if err != nil {
		return errors.Handle(c, err)
	}
	
	return c.JSON(fiber.Map{
		"status": "success",
		"message": "Предмет успешно приобретен",
		"profile": profile,
	})
}

// GetUserWords получает словарь пользователя
func (h *UserHandler) GetUserWords(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	
	words, err := h.Service.GetUserWords(userID)
	if err != nil {
		return errors.Handle(c, err)
	}
	
	return c.JSON(fiber.Map{
		"words": words,
	})
}

// LearnWord добавляет слово в словарь пользователя
func (h *UserHandler) LearnWord(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	
	var req dto.LearnWordRequest
	if err := utils.ParseAndValidate(c, &req); err != nil {
		return errors.Handle(c, err)
	}
	
	if err := h.Service.LearnWord(userID, &req); err != nil {
		return errors.Handle(c, err)
	}
	
	return c.JSON(fiber.Map{
		"status": "success",
		"message": "Слово успешно добавлено в словарь",
	})
}
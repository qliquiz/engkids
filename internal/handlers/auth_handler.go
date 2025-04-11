package handlers

import (
	"engkids/internal/dto"
	"engkids/internal/errors"
	"engkids/internal/services"
	"engkids/pkg/utils"

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

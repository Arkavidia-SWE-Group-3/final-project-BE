package handlers

import (
	"Go-Starter-Template/pkg/company"

	"github.com/gofiber/fiber/v2"

	"github.com/go-playground/validator/v10"

	"Go-Starter-Template/domain"

	"Go-Starter-Template/internal/api/presenters"
)

type (
	CompanyHandler interface {
		GetProfile(c *fiber.Ctx) error
	}
	companyHandler struct {
		CompanyService company.CompanyService
		Validator      *validator.Validate
	}
)

func NewCompanyHandler(companyService company.CompanyService, validator *validator.Validate) CompanyHandler {
	return &companyHandler{
		CompanyService: companyService,
		Validator:      validator,
	}
}

func (h *companyHandler) GetProfile(c *fiber.Ctx) error {
	slug := c.Params("slug")

	res, err := h.CompanyService.GetProfile(c.Context(), slug)

	if err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedGetProfile, err)
	}

	return presenters.SuccessResponse(c, res, fiber.StatusOK, domain.MessageSuccessGetProfile)
}

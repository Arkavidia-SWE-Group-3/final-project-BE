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
		AddJob(c *fiber.Ctx) error
		UpdateJob(c *fiber.Ctx) error
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

func (h *companyHandler) AddJob(c *fiber.Ctx) error {
	var req domain.CompanyAddJobRequest

	if err := c.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedAddJob, err)
	}

	if err := h.Validator.Struct(req); err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedAddJob, err)
	}

	err := h.CompanyService.AddJob(c.Context(), req.CompanyID, req)

	if err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedAddJob, err)
	}

	return presenters.SuccessResponse(c, nil, fiber.StatusCreated, domain.MessageSuccessAddJob)
}

func (h *companyHandler) UpdateJob(c *fiber.Ctx) error {
	var req domain.CompanyUpdateJobRequest

	if err := c.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedAddJob, err)
	}

	if err := h.Validator.Struct(req); err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedAddJob, err)
	}

	err := h.CompanyService.UpdateJob(c.Context(), req.CompanyID, req)

	if err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedAddJob, err)
	}

	return presenters.SuccessResponse(c, nil, fiber.StatusCreated, domain.MessageSuccessAddJob)
}

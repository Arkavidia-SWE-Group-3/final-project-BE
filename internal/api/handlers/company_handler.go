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
		UpdateProfile(c *fiber.Ctx) error
		RegisterCompany(c *fiber.Ctx) error
		LoginCompany(c *fiber.Ctx) error
		GetListCompany(c *fiber.Ctx) error
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

func (h *companyHandler) GetListCompany(c *fiber.Ctx) error {
	res, err := h.CompanyService.GetListCompany(c.Context())

	if err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedGetListCompany, err)
	}

	return presenters.SuccessResponse(c, res, fiber.StatusOK, domain.MessageSuccessGetListCompany)
}

func (h *companyHandler) LoginCompany(c *fiber.Ctx) error {
	var req domain.CompanyLoginRequest

	if err := c.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedLogin, err)
	}

	if err := h.Validator.Struct(req); err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedLogin, err)
	}

	res, err := h.CompanyService.LoginCompany(c.Context(), req)

	if err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedLogin, err)
	}

	return presenters.SuccessResponse(c, res, fiber.StatusOK, domain.MessageSuccessLogin)
}

func (h *companyHandler) RegisterCompany(c *fiber.Ctx) error {
	var req domain.CompanyRegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedRegisterCompany, err)
	}

	if err := h.Validator.Struct(req); err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedRegisterCompany, err)
	}

	err := h.CompanyService.RegisterCompany(c.Context(), req)

	if err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedRegisterCompany, err)
	}

	return presenters.SuccessResponse(c, nil, fiber.StatusCreated, domain.MessageSuccessRegisterCompany)
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

	userID := c.Locals("user_id").(string)

	err := h.CompanyService.AddJob(c.Context(), req, userID)

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

	userID := c.Locals("user_id").(string)

	err := h.CompanyService.UpdateJob(c.Context(), req, userID)

	if err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedAddJob, err)
	}

	return presenters.SuccessResponse(c, nil, fiber.StatusCreated, domain.MessageSuccessAddJob)
}

func (h *companyHandler) UpdateProfile(c *fiber.Ctx) error {
	var req domain.CompanyUpdateProfileRequest

	if err := c.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedUpdateProfileCompany, err)
	}

	if err := h.Validator.Struct(req); err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedUpdateProfileCompany, err)
	}

	userID := c.Locals("user_id").(string)

	req.Logo, _ = c.FormFile("logo")
	req.Headline, _ = c.FormFile("cover")

	err := h.CompanyService.UpdateProfile(c.Context(), req, userID)

	if err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedUpdateProfileCompany, err)
	}

	return presenters.SuccessResponse(c, nil, fiber.StatusCreated, domain.MessageSuccessUpdateProfileCompany)
}

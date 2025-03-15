package handlers

import (
	"Go-Starter-Template/domain"
	"Go-Starter-Template/pkg/job"

	"Go-Starter-Template/internal/api/presenters"

	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/go-playground/validator/v10"
)

type (
	JobHandler interface {
		SearchJob(c *fiber.Ctx) error
		GetJobDetail(c *fiber.Ctx) error
	}
	jobHandler struct {
		JobService job.JobService
		Validator  *validator.Validate
	}
)

func NewJobHandler(jobService job.JobService, validator *validator.Validate) JobHandler {
	return &jobHandler{
		JobService: jobService,
		Validator:  validator,
	}
}

func (h *jobHandler) GetJobDetail(c *fiber.Ctx) error {
	id := c.Params("id")

	res, err := h.JobService.GetJobDetail(c.Context(), id)

	if err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedGetJobDetail, err)
	}

	return presenters.SuccessResponse(c, res, fiber.StatusOK, domain.MessageSuccessGetJobDetail)
}

func (h *jobHandler) SearchJob(c *fiber.Ctx) error {

	var title = c.Query("title")
	var jobType = c.Query("job_type")
	var locationType = c.Query("location_type")
	var minSalaryStr = c.Query("min_salary")
	var maxSalaryStr = c.Query("max_salary")
	var experienceLevel = c.Query("experience_level")
	var sortBy = c.Query("sort_by")

	minSalary, err := strconv.Atoi(minSalaryStr)

	if err != nil {
		minSalary = 0
	}

	maxSalary, err := strconv.Atoi(maxSalaryStr)

	if err != nil {
		maxSalary = 0
	}

	var datePosted = c.Query("date_posted")

	var jobSearchRequest = domain.JobSearchRequest{
		Title:           title,
		JobType:         jobType,
		LocationType:    locationType,
		ExperienceLevel: experienceLevel,
		MinSalary:       minSalary,
		MaxSalary:       maxSalary,
		SortBy:          sortBy,
		DatePosted:      datePosted,
	}

	res, err := h.JobService.SearchJob(c.Context(), jobSearchRequest)

	if err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedGetJobs, err)
	}

	return presenters.SuccessResponse(c, res, fiber.StatusOK, domain.MessageSuccessSearchJobs)
}

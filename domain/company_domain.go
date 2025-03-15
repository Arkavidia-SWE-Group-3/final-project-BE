package domain

import (
	"errors"
	"mime/multipart"
)

var (
	MessageSuccessAddJob               = "Job added successfully"
	MessageSuccessUpdateProfileCompany = "Company profile updated successfully"

	MessageFailedAddJob               = "Failed to add job"
	MessageFailedUpdateProfileCompany = "Failed to update company profile"

	ErrJobNotCreated     = errors.New("job not created")
	ErrJobNotUpdated     = errors.New("job not updated")
	ErrCompanyNotUpdated = errors.New("company not updated")
)

type (
	CompanyProfileResponse struct {
		CompanyInfo CompanyInfoResponse   `json:"company_info"`
		ComapnyJobs []CompanyJobsResponse `json:"company_jobs"`
	}

	CompanyInfoResponse struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		About    string `json:"about"`
		Industry string `json:"industry"`
		Logo     string `json:"logo"`
		Headline string `json:"cover"`
	}

	CompanyJobsResponse struct {
		ID              string                     `json:"id"`
		Title           string                     `json:"title"`
		Location        string                     `json:"location"`
		LocationType    string                     `json:"location_type"`
		JobType         string                     `json:"job_type"`
		ExperienceLevel string                     `json:"experience"`
		SalaryMin       int                        `json:"min_salary"`
		SalaryMax       int                        `json:"max_salary"`
		Status          string                     `json:"status"`
		Description     string                     `json:"description"`
		Skills          []CompanyJobSkillsResponse `json:"skills"`
		Posted          string                     `json:"posted"`
	}

	CompanyJobSkillsResponse struct {
		ID      string `json:"id"`
		SkillID string `json:"skill_id"`
		Name    string `json:"name"`
	}

	CompanyAddJobRequest struct {
		CompanyID       string `json:"company_id"`
		Title           string `json:"title"`
		Location        string `json:"location"`
		LocationType    string `json:"location_type"`
		JobType         string `json:"job_type"`
		ExperienceLevel string `json:"experience"`
		SalaryMin       int    `json:"min_salary"`
		SalaryMax       int    `json:"max_salary"`
		Description     string `json:"description"`
		Skills          []string
		Status          string `json:"status"`
	}

	CompanyUpdateJobRequest struct {
		CompanyID       string `json:"company_id"`
		JobID           string `json:"job_id"`
		Title           string `json:"title"`
		Location        string `json:"location"`
		LocationType    string `json:"location_type"`
		JobType         string `json:"job_type"`
		ExperienceLevel string `json:"experience"`
		SalaryMin       int    `json:"min_salary"`
		SalaryMax       int    `json:"max_salary"`
		Description     string `json:"description"`
		Skills          []string
	}

	CompanyUpdateProfileRequest struct {
		CompanyID string                `json:"company_id" form:"company_id" validate:"required"`
		Name      string                `json:"name"`
		Industry  string                `json:"industry"`
		Logo      *multipart.FileHeader `json:"logo"`
		Headline  *multipart.FileHeader `json:"cover"`
	}
)

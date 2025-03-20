package domain

import (
	"errors"
	"mime/multipart"
)

var (
	MessageSuccessAddJob               = "Job added successfully"
	MessageSuccessUpdateProfileCompany = "Company profile updated successfully"
	MessageSuccessRegisterCompany      = "Company registered successfully"
	MessageSuccessGetListCompany       = "Company profile retrieved successfully"

	MessageFailedAddJob               = "Failed to add job"
	MessageFailedUpdateProfileCompany = "Failed to update company profile"
	MessageFailedRegisterCompany      = "Failed to register company"
	MessageFailedGetListCompany       = "Failed to retrieve company profile"

	ErrJobNotCreated            = errors.New("job not created")
	ErrJobNotUpdated            = errors.New("job not updated")
	ErrCompanyNotUpdated        = errors.New("company not updated")
	ErrCompanyNotRegistered     = errors.New("company not created")
	ErrCompanyAlreadyRegistered = errors.New("company already registered")
	ErrCompanyNotFound          = errors.New("company not found")
)

type (
	CompanyProfileResponse struct {
		CompanyInfo  CompanyInfoResponse    `json:"company_info"`
		ComapnyJobs  []CompanyJobsResponse  `json:"company_jobs"`
		CompanyPosts []CompanyPostsResponse `json:"company_posts"`
	}

	CompanyPostsResponse struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		ProfilePicture string `json:"profile_picture"`
		Content        string `json:"content"`
		CreatedAt      string `json:"created_at"`
		Headline       string `json:"headline"`
		Asset          string `json:"asset"`
		Slug           string `json:"slug"`
		Type           string `json:"type"`
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

	CompanyLoginResponse struct {
		Email          string `json:"email"`
		Token          string `json:"token"`
		Role           string `json:"role"`
		Slug           string `json:"slug"`
		Name           string `json:"name"`
		CurrentTitle   string `json:"current_title"`
		ProfilePicture string `json:"profile_picture"`
	}

	CompanyJobSkillsResponse struct {
		ID      string `json:"id"`
		SkillID string `json:"skill_id"`
		Name    string `json:"name"`
	}

	CompanyListResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
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

	CompanyRegisterRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		About    string `json:"about"`
		Industry string `json:"industry"`
	}

	CompanyLoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
		Name      string                `json:"name" form:"name"`
		Industry  string                `json:"industry" form:"industry"`
		Logo      *multipart.FileHeader `json:"logo" form:"logo"`
		Headline  *multipart.FileHeader `json:"cover" form:"cover"`
		About     string                `json:"about" form:"about"`
	}
)

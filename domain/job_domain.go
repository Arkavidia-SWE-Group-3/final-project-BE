package domain

import "mime/multipart"

var (
	MessageFailedGetJobs      = "Failed to get jobs"
	MessageFailedSearchJobs   = "Failed to search jobs"
	MessageFailedGetJobDetail = "Successfully get job detail"
	MessageFailedApplyJob     = "Failed to apply job"

	MessageSuccessSearchJobs   = "Successfully search jobs"
	MessageSuccessGetJobDetail = "Successfully get job detail"
	MessageSuccessApplyJob     = "Successfully apply job"
)

type (
	JobSearchRequest struct {
		Title           string `json:"title"`
		JobType         string `json:"job_type"`
		LocationType    string `json:"location_type"`
		ExperienceLevel string `json:"experience_level"`
		MinSalary       int    `json:"min_salary"`
		MaxSalary       int    `json:"max_salary"`
		DatePosted      string `json:"date_posted"`
		SortBy          string `json:"sort_by"`
	}

	JobApplyRequest struct {
		JobID  string                `json:"job_id" form:"job_id"`
		Resume *multipart.FileHeader `json:"resume" form:"resume"`
	}

	JobSearchResponse struct {
		ID              string   `json:"id"`
		CompanyName     string   `json:"company"`
		CompanyLogo     string   `json:"logo"`
		CompanySlug     string   `json:"company_slug"`
		Title           string   `json:"title"`
		Location        string   `json:"location"`
		LocationType    string   `json:"location_type"`
		JobType         string   `json:"type"`
		ExperienceLevel string   `json:"experience"`
		SalaryMin       int      `json:"min_salary"`
		SalaryMax       int      `json:"max_salary"`
		Description     string   `json:"description"`
		Status          string   `json:"status"`
		Posted          string   `json:"posted"`
		Skills          []string `json:"skills"`
	}

	JobDetailResponse struct {
		ID              string   `json:"id"`
		CompanyName     string   `json:"company"`
		CompanyLogo     string   `json:"logo"`
		CompanySlug     string   `json:"company_slug"`
		Title           string   `json:"title"`
		Location        string   `json:"location"`
		LocationType    string   `json:"location_type"`
		JobType         string   `json:"type"`
		ExperienceLevel string   `json:"experience"`
		SalaryMin       int      `json:"min_salary"`
		SalaryMax       int      `json:"max_salary"`
		Description     string   `json:"description"`
		Status          string   `json:"status"`
		Posted          string   `json:"posted"`
		Skills          []string `json:"skills"`
	}
)

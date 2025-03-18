package domain

import "mime/multipart"

var (
	MessageFailedGetJobs                 = "Failed to get jobs"
	MessageFailedSearchJobs              = "Failed to search jobs"
	MessageFailedGetJobDetail            = "Successfully get job detail"
	MessageFailedApplyJob                = "Failed to apply job"
	MessageFailedGetApplicants           = "Failed to get applicants"
	MessageFailedChangeApplicationStatus = "Failed to change application status"

	MessageSuccessSearchJobs              = "Successfully search jobs"
	MessageSuccessGetJobDetail            = "Successfully get job detail"
	MessageSuccessApplyJob                = "Successfully apply job"
	MessageSuccessGetApplicants           = "Successfully get applicants"
	MessageSuccessChangeApplicationStatus = "Successfully change application status"
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

	JobChangeApplicationStatusRequest struct {
		JobApplicationID  string `json:"applicant_id"`
		ApplicationStatus string `json:"status"`
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

	JobApplicantResponse struct {
		ID                 string `json:"id"`
		UserID             string `json:"user_id"`
		UserName           string `json:"name"`
		UserSlug           string `json:"slug"`
		UserProfilePicture string `json:"profile_picture"`
		UserHeadline       string `json:"headline"`
		ResumeURL          string `json:"resume_url"`
		Status             string `json:"status"`
		AppliedAt          string `json:"applied_at"`
	}
)

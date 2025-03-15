package domain

var (
	MessageFailedGetJobs = "Failed to get jobs"

	MessageSuccessSearchJobs = "Successfully search jobs"
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

	JobSearchResponse struct {
		ID              string   `json:"id"`
		CompanyName     string   `json:"company"`
		CompanyLogo     string   `json:"logo"`
		CompanySlug     string   `json:"slug"`
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

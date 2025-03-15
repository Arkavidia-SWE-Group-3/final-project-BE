package domain

var ()

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
	}

	CompanyJobsResponse struct {
		ID              string                     `json:"id"`
		Title           string                     `json:"title"`
		Location        string                     `json:"location"`
		LocationType    string                     `json:"location_type"`
		JobType         string                     `json:"job_type"`
		ExperienceLevel string                     `json:"experience_level"`
		SalaryMin       int                        `json:"salary_min"`
		SalaryMax       int                        `json:"salary_max"`
		Status          string                     `json:"status"`
		Description     string                     `json:"description"`
		Skills          []CompanyJobSkillsResponse `json:"skills"`
	}

	CompanyJobSkillsResponse struct {
		ID      string `json:"id"`
		SkillID string `json:"skill_id"`
		Name    string `json:"name"`
	}
)

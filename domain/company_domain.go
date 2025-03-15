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
)

package job

import (
	"Go-Starter-Template/domain"
	"Go-Starter-Template/internal/utils"
	"Go-Starter-Template/internal/utils/storage"
	jwtService "Go-Starter-Template/pkg/jwt"
	"context"
)

type (
	JobService interface {
		SearchJob(ctx context.Context, jobFilters domain.JobSearchRequest) ([]domain.JobSearchResponse, error)
		GetJobDetail(ctx context.Context, id string) (domain.JobDetailResponse, error)
	}

	jobService struct {
		jobRepository JobRepository
		awsS3         storage.AwsS3
		jwtService    jwtService.JWTService
	}
)

func NewJobService(jobRepository JobRepository, awsS3 storage.AwsS3, jwtService jwtService.JWTService) JobService {
	return &jobService{jobRepository: jobRepository, awsS3: awsS3, jwtService: jwtService}
}

func (s *jobService) GetJobDetail(ctx context.Context, id string) (domain.JobDetailResponse, error) {
	res, err := s.jobRepository.GetJobDetail(ctx, id)

	var jobResult domain.JobDetailResponse

	var jobSkills []string

	for _, skill := range res.Skills {
		jobSkills = append(jobSkills, skill.Name)
	}

	if jobSkills == nil {
		jobSkills = []string{}
	}

	jobResult = domain.JobDetailResponse{
		ID:              res.ID.String(),
		CompanyName:     res.Company.Name,
		CompanySlug:     res.Company.Slug,
		CompanyLogo:     "",
		Title:           res.Title,
		Location:        res.Location,
		LocationType:    res.LocationType,
		JobType:         res.JobType,
		ExperienceLevel: res.ExperienceLevel,
		SalaryMin:       res.SalaryMin,
		SalaryMax:       res.SalaryMax,
		Description:     res.Description,
		Status:          res.Status,
		Posted:          utils.ConvertTimeToString(res.CreatedAt),
		Skills:          jobSkills,
	}

	if err != nil {
		return domain.JobDetailResponse{}, err
	}

	return jobResult, nil
}

func (s *jobService) SearchJob(ctx context.Context, jobFilters domain.JobSearchRequest) ([]domain.JobSearchResponse, error) {
	res, err := s.jobRepository.SearchJob(ctx, jobFilters)

	if err != nil {
		return nil, err
	}

	var jobSearchResponse []domain.JobSearchResponse

	for _, job := range res {
		var jobSkills []string

		for _, skill := range job.Skills {
			jobSkills = append(jobSkills, skill.Name)
		}

		if jobSkills == nil {
			jobSkills = []string{}
		}

		jobSearchResponse = append(jobSearchResponse, domain.JobSearchResponse{
			ID:              job.ID.String(),
			CompanyName:     job.Company.Name,
			CompanySlug:     job.Company.Slug,
			CompanyLogo:     "",
			Title:           job.Title,
			Location:        job.Location,
			LocationType:    job.LocationType,
			JobType:         job.JobType,
			ExperienceLevel: job.ExperienceLevel,
			SalaryMin:       job.SalaryMin,
			SalaryMax:       job.SalaryMax,
			Description:     job.Description,
			Status:          job.Status,
			Posted:          utils.ConvertTimeToString(job.CreatedAt),
			Skills:          jobSkills,
		})
	}

	if jobSearchResponse == nil {
		jobSearchResponse = []domain.JobSearchResponse{}
	}

	return jobSearchResponse, nil
}

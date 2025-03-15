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

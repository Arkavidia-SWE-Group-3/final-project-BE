package company

import (
	"Go-Starter-Template/domain"
	"Go-Starter-Template/internal/utils"
	"Go-Starter-Template/internal/utils/storage"
	jwtService "Go-Starter-Template/pkg/jwt"
	"context"
)

type (
	CompanyService interface {
		GetProfile(ctx context.Context, slug string) (*domain.CompanyProfileResponse, error)
	}

	companyService struct {
		companyRepository CompanyRepository
		awsS3             storage.AwsS3
		jwtService        jwtService.JWTService
	}
)

func NewCompanyService(companyRepository CompanyRepository, awsS3 storage.AwsS3, jwtService jwtService.JWTService) CompanyService {
	return &companyService{companyRepository: companyRepository, awsS3: awsS3, jwtService: jwtService}
}

func (s *companyService) GetProfile(ctx context.Context, slug string) (*domain.CompanyProfileResponse, error) {
	company, err := s.companyRepository.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	companyJobs, err := s.companyRepository.GetJobsByCompanyID(ctx, company.ID)
	if err != nil {
		return nil, err
	}

	companyInfoResponse := domain.CompanyInfoResponse{
		ID:       company.ID.String(),
		Name:     company.Name,
		About:    company.About,
		Industry: company.Industry,
		Logo:     "",
		Headline: "",
	}

	var companyJobsResponse []domain.CompanyJobsResponse
	for _, job := range companyJobs {
		skills, err := s.companyRepository.GetJobSkillsByJobID(ctx, job.ID)
		if err != nil {
			return nil, err
		}

		var companyJobSkillsResponse []domain.CompanyJobSkillsResponse
		for _, skill := range skills {
			companyJobSkillsResponse = append(companyJobSkillsResponse, domain.CompanyJobSkillsResponse{
				ID:      skill.ID.String(),
				SkillID: skill.SkillID.String(),
				Name:    skill.Skill.Name,
			})
		}

		if companyJobSkillsResponse == nil {
			companyJobSkillsResponse = []domain.CompanyJobSkillsResponse{}
		}

		companyJobsResponse = append(companyJobsResponse, domain.CompanyJobsResponse{
			ID:              job.ID.String(),
			Title:           job.Title,
			Location:        job.Location,
			LocationType:    job.LocationType,
			JobType:         job.JobType,
			ExperienceLevel: job.ExperienceLevel,
			SalaryMin:       job.SalaryMin,
			SalaryMax:       job.SalaryMax,
			Status:          job.Status,
			Description:     job.Description,
			Skills:          companyJobSkillsResponse,
			Posted:          utils.ConvertTimeToString(job.CreatedAt),
		})
	}

	if companyJobsResponse == nil {
		companyJobsResponse = []domain.CompanyJobsResponse{}

	}

	return &domain.CompanyProfileResponse{
		CompanyInfo: companyInfoResponse,
		ComapnyJobs: companyJobsResponse,
	}, nil
}

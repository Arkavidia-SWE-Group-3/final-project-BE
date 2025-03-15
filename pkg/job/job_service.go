package job

import (
	"Go-Starter-Template/domain"
	"Go-Starter-Template/entities"
	"Go-Starter-Template/internal/utils"
	"Go-Starter-Template/internal/utils/storage"
	jwtService "Go-Starter-Template/pkg/jwt"
	"context"

	"github.com/google/uuid"
)

type (
	JobService interface {
		SearchJob(ctx context.Context, jobFilters domain.JobSearchRequest) ([]domain.JobSearchResponse, error)
		GetJobDetail(ctx context.Context, id string) (domain.JobDetailResponse, error)
		ApplyJob(ctx context.Context, req domain.JobApplyRequest, userID string) error
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

func (s *jobService) ApplyJob(ctx context.Context, req domain.JobApplyRequest, userID string) error {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return domain.ErrParseUUID
	}

	parsedJobID, err := uuid.Parse(req.JobID)

	if err != nil {
		return domain.ErrParseUUID
	}

	jobApplication := entities.JobApplication{
		ID:     uuid.New(),
		JobID:  parsedJobID,
		UserID: parsedUserID,
		Status: "Under Review",
	}

	allowedMimeTypes := []string{"application/pdf"}

	if req.Resume != nil {
		objectKey, err := s.awsS3.UploadFile(userID, req.Resume, "resume", allowedMimeTypes...)

		if err != nil {
			return domain.ErrUploadFile
		}

		jobApplication.CV = s.awsS3.GetPublicLinkKey(objectKey)

		err = s.jobRepository.ApplyJob(ctx, jobApplication)

		if err != nil {
			return err
		}
	}

	return nil
}

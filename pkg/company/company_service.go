package company

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
	CompanyService interface {
		GetProfile(ctx context.Context, slug string) (*domain.CompanyProfileResponse, error)
		AddJob(ctx context.Context, req domain.CompanyAddJobRequest, userID string) error
		UpdateJob(ctx context.Context, req domain.CompanyUpdateJobRequest, userID string) error
		UpdateProfile(ctx context.Context, req domain.CompanyUpdateProfileRequest, userID string) error
		LoginCompany(ctx context.Context, req domain.CompanyLoginRequest) (*domain.CompanyLoginResponse, error)
		RegisterCompany(ctx context.Context, req domain.CompanyRegisterRequest) error
		GetListCompany(ctx context.Context) ([]domain.CompanyListResponse, error)
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

func (s *companyService) GetListCompany(ctx context.Context) ([]domain.CompanyListResponse, error) {
	companies, err := s.companyRepository.GetListCompany(ctx)

	if err != nil {
		return nil, domain.ErrCompanyNotFound
	}

	var companyListResponse []domain.CompanyListResponse

	for _, company := range companies {
		companyListResponse = append(companyListResponse, domain.CompanyListResponse{
			ID:   company.ID.String(),
			Name: company.Name,
		})

	}

	return companyListResponse, nil
}

func (s *companyService) LoginCompany(ctx context.Context, req domain.CompanyLoginRequest) (*domain.CompanyLoginResponse, error) {
	user, company, err := s.companyRepository.GetCompanyByEmail(ctx, req.Email)

	if err != nil {
		return nil, domain.ErrCompanyNotFound
	}

	if user.Role != "company" {
		return nil, domain.ErrCompanyNotFound
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, domain.CredentialInvalid
	}

	token := s.jwtService.GenerateTokenUser(user.ID.String(), user.Role)

	return &domain.CompanyLoginResponse{
		Email:          user.Email,
		Token:          token,
		Role:           user.Role,
		Slug:           company.Slug,
		Name:           company.Name,
		CurrentTitle:   "",
		ProfilePicture: user.ProfilePicture,
	}, nil
}

func (s *companyService) RegisterCompany(ctx context.Context, req domain.CompanyRegisterRequest) error {
	_, _, err := s.companyRepository.GetCompanyByEmail(ctx, req.Email)

	if err == nil {
		return domain.ErrCompanyAlreadyRegistered
	}

	password, err := utils.HashPassword(req.Password)

	if err != nil {
		return err
	}

	validated := utils.ValidatePassword(req.Password)

	if !validated {
		return domain.ErrPasswordNotValid
	}

	company := entities.Companies{
		Name:     req.Name,
		Slug:     utils.CreateSlug(req.Name),
		About:    req.About,
		Industry: req.Industry,
	}

	user := entities.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: password,
		Role:     "company",
		Slug:     utils.CreateSlug(req.Name),
	}

	err = s.companyRepository.RegisterCompany(ctx, company, user)

	if err != nil {
		return domain.ErrCompanyNotRegistered
	}

	return nil
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

	companyPosts, err := s.companyRepository.GetPostsByCompanyID(ctx, company.UserID)

	if err != nil {
		return nil, err
	}

	companyInfoResponse := domain.CompanyInfoResponse{
		ID:       company.ID.String(),
		Name:     company.Name,
		About:    company.About,
		Industry: company.Industry,
		Logo:     company.User.ProfilePicture,
		Headline: company.User.Headline,
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

	var companyPostsResponse []domain.CompanyPostsResponse
	for _, post := range companyPosts {
		companyPostsResponse = append(companyPostsResponse, domain.CompanyPostsResponse{
			ID:             post.ID.String(),
			Name:           post.User.Name,
			ProfilePicture: post.User.ProfilePicture,
			Content:        post.Content,
			CreatedAt:      utils.ConvertTimeToString(post.CreatedAt),
			Headline:       post.User.CurrentTitle,
			Asset:          post.Asset,
			Slug:           company.Slug,
			Type:           post.User.Role,
		})
	}

	if companyPostsResponse == nil {
		companyPostsResponse = []domain.CompanyPostsResponse{}
	}

	return &domain.CompanyProfileResponse{
		CompanyInfo:  companyInfoResponse,
		ComapnyJobs:  companyJobsResponse,
		CompanyPosts: companyPostsResponse,
	}, nil
}

func (s *companyService) UpdateProfile(ctx context.Context, req domain.CompanyUpdateProfileRequest, userID string) error {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return domain.ErrParseUUID
	}

	company_id, err := s.companyRepository.GetCompanyByUserID(ctx, parsedUserID)

	if err != nil {
		return domain.ErrCompanyNotFound
	}

	company := entities.Companies{
		ID:       company_id.ID,
		Name:     req.Name,
		Industry: req.Industry,
		About:    req.About,
	}

	user := entities.User{
		ID: parsedUserID,
	}

	allowedMimetype := []string{"image/jpeg", "image/png"}

	if req.Logo != nil {
		objectKey, err := s.awsS3.UploadFile(utils.GenerateRandomFileName(req.Logo.Filename), req.Logo, "profile-picture", allowedMimetype...)

		if err != nil {
			return domain.ErrUploadFile
		}

		user.ProfilePicture = s.awsS3.GetPublicLinkKey(objectKey)
	}

	if req.Headline != nil {
		objectKey, err := s.awsS3.UploadFile(utils.GenerateRandomFileName(req.Headline.Filename), req.Headline, "headline", allowedMimetype...)

		if err != nil {
			return domain.ErrUploadFile
		}

		user.Headline = s.awsS3.GetPublicLinkKey(objectKey)
	}

	err = s.companyRepository.UpdateProfile(ctx, company, user)

	if err != nil {
		return domain.ErrCompanyNotUpdated
	}

	return nil
}

func (s *companyService) AddJob(ctx context.Context, req domain.CompanyAddJobRequest, userID string) error {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return domain.ErrParseUUID
	}

	company, err := s.companyRepository.GetCompanyByUserID(ctx, parsedUserID)

	if err != nil {
		return domain.ErrCompanyNotFound
	}

	job := entities.Job{
		ID:              uuid.New(),
		CompanyID:       company.ID,
		Title:           req.Title,
		Location:        req.Location,
		LocationType:    req.LocationType,
		JobType:         req.JobType,
		ExperienceLevel: req.ExperienceLevel,
		SalaryMin:       req.SalaryMin,
		SalaryMax:       req.SalaryMax,
		Description:     req.Description,
		Status:          "active",
	}

	jobID := s.companyRepository.AddJob(ctx, job)

	if jobID == uuid.Nil {
		return domain.ErrJobNotCreated
	}

	for _, skillID := range req.Skills {
		jobSkill := entities.JobSkill{
			JobID:   jobID,
			SkillID: uuid.MustParse(skillID),
		}

		err := s.companyRepository.AddJobSkill(ctx, jobSkill)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *companyService) UpdateJob(ctx context.Context, req domain.CompanyUpdateJobRequest, userID string) error {

	parsedUserID, err_parsed := uuid.Parse(userID)

	if err_parsed != nil {
		return domain.ErrParseUUID
	}

	companyID, err_get_company_id := s.companyRepository.GetCompanyByUserID(ctx, parsedUserID)

	if err_get_company_id != nil {
		return domain.ErrCompanyNotFound
	}

	job := entities.Job{
		ID:              uuid.MustParse(req.JobID),
		CompanyID:       companyID.ID,
		Title:           req.Title,
		Location:        req.Location,
		LocationType:    req.LocationType,
		JobType:         req.JobType,
		ExperienceLevel: req.ExperienceLevel,
		SalaryMin:       req.SalaryMin,
		SalaryMax:       req.SalaryMax,
		Description:     req.Description,
		Status:          "active",
	}

	err := s.companyRepository.UpdateJob(ctx, job)

	if err != nil {
		return domain.ErrJobNotUpdated
	}

	err = s.companyRepository.DeleteJobSkillsByJobID(ctx, job.ID)

	if err != nil {
		return err
	}

	for _, skillID := range req.Skills {
		jobSkill := entities.JobSkill{
			JobID:   job.ID,
			SkillID: uuid.MustParse(skillID),
		}

		err := s.companyRepository.AddJobSkill(ctx, jobSkill)
		if err != nil {
			return err
		}
	}

	return nil
}

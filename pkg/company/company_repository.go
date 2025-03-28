package company

import (
	"Go-Starter-Template/entities"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	CompanyRepository interface {
		GetBySlug(ctx context.Context, slug string) (entities.Companies, error)
		GetJobsByCompanyID(ctx context.Context, companyID uuid.UUID) ([]entities.Job, error)
		GetJobSkillsByJobID(ctx context.Context, jobID uuid.UUID) ([]entities.JobSkill, error)
		AddJob(ctx context.Context, job entities.Job) uuid.UUID
		AddJobSkill(ctx context.Context, jobSkill entities.JobSkill) error
		UpdateJob(ctx context.Context, job entities.Job) error
		DeleteJobSkillsByJobID(ctx context.Context, jobID uuid.UUID) error
		UpdateProfile(ctx context.Context, company entities.Companies, user entities.User) error
		RegisterCompany(ctx context.Context, company entities.Companies, user entities.User) error
		GetCompanyByEmail(ctx context.Context, email string) (entities.User, entities.Companies, error)
		GetCompanyByUserID(ctx context.Context, userID uuid.UUID) (entities.Companies, error)
		GetPostsByCompanyID(ctx context.Context, companyID uuid.UUID) ([]entities.Post, error)
		GetListCompany(ctx context.Context) ([]entities.Companies, error)
	}
	companyRepository struct {
		db *gorm.DB
	}
)

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) GetListCompany(ctx context.Context) ([]entities.Companies, error) {
	var companies []entities.Companies
	if err := r.db.WithContext(ctx).Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}

func (r *companyRepository) GetCompanyByUserID(ctx context.Context, userID uuid.UUID) (entities.Companies, error) {
	var company entities.Companies
	if err := r.db.WithContext(ctx).Preload("User").First(&company, "user_id = ?", userID).Error; err != nil {
		return entities.Companies{}, err
	}
	return company, nil
}

func (r *companyRepository) GetCompanyByEmail(ctx context.Context, email string) (entities.User, entities.Companies, error) {
	var user entities.User
	var company entities.Companies

	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return entities.User{}, entities.Companies{}, err
	}

	if err := r.db.WithContext(ctx).First(&company, "user_id = ?", user.ID).Error; err != nil {
		return entities.User{}, entities.Companies{}, err
	}

	return user, company, nil
}

func (r *companyRepository) RegisterCompany(ctx context.Context, company entities.Companies, user entities.User) error {
	companyID := uuid.New()
	userID := uuid.New()

	company.ID = companyID
	user.ID = userID

	company.UserID = userID

	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(&company).Error; err != nil {
		return err
	}

	return nil
}

func (r *companyRepository) GetBySlug(ctx context.Context, slug string) (entities.Companies, error) {
	var company entities.Companies
	if err := r.db.WithContext(ctx).Preload("User").First(&company, "slug = ?", slug).Error; err != nil {
		return entities.Companies{}, err
	}

	return company, nil

}

func (r *companyRepository) UpdateProfile(ctx context.Context, company entities.Companies, user entities.User) error {
	if err := r.db.WithContext(ctx).Model(&company).Updates(&company).Error; err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Model(&user).Updates(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *companyRepository) GetJobsByCompanyID(ctx context.Context, companyID uuid.UUID) ([]entities.Job, error) {
	var jobs []entities.Job

	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *companyRepository) GetPostsByCompanyID(ctx context.Context, companyID uuid.UUID) ([]entities.Post, error) {
	var posts []entities.Post

	if err := r.db.WithContext(ctx).Preload("User").Order("created_at desc").Where("user_id = ?", companyID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *companyRepository) GetJobSkillsByJobID(ctx context.Context, jobID uuid.UUID) ([]entities.JobSkill, error) {
	var jobSkill []entities.JobSkill

	if err := r.db.Preload("Skill").WithContext(ctx).Where("job_id = ?", jobID).Find(&jobSkill).Error; err != nil {
		return nil, err
	}

	return jobSkill, nil
}

func (r *companyRepository) AddJob(ctx context.Context, job entities.Job) uuid.UUID {
	var jobID uuid.UUID

	if err := r.db.WithContext(ctx).Create(&job).Error; err != nil {
		return jobID
	}

	jobID = job.ID
	return jobID

}

func (r *companyRepository) AddJobSkill(ctx context.Context, jobSkill entities.JobSkill) error {
	if err := r.db.WithContext(ctx).Create(&jobSkill).Error; err != nil {
		return err
	}
	return nil
}

func (r *companyRepository) UpdateJob(ctx context.Context, job entities.Job) error {
	if err := r.db.WithContext(ctx).Model(&job).Updates(&job).Error; err != nil {
		return err
	}
	return nil
}

func (r *companyRepository) DeleteJobSkillsByJobID(ctx context.Context, jobID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("job_id = ?", jobID).Delete(&entities.JobSkill{}).Error; err != nil {
		return err
	}
	return nil
}

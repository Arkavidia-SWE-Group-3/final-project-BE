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
	}
	companyRepository struct {
		db *gorm.DB
	}
)

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) GetBySlug(ctx context.Context, slug string) (entities.Companies, error) {
	var company entities.Companies
	if err := r.db.WithContext(ctx).First(&company, "slug = ?", slug).Error; err != nil {
		return entities.Companies{}, err
	}
	return company, nil

}

func (r *companyRepository) GetJobsByCompanyID(ctx context.Context, companyID uuid.UUID) ([]entities.Job, error) {
	var jobs []entities.Job

	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *companyRepository) GetJobSkillsByJobID(ctx context.Context, jobID uuid.UUID) ([]entities.JobSkill, error) {
	var jobSkill []entities.JobSkill

	if err := r.db.Preload("Skill").WithContext(ctx).Where("job_id = ?", jobID).Find(&jobSkill).Error; err != nil {
		return nil, err
	}

	return jobSkill, nil
}

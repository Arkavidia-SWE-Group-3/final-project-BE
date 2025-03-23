package job

import (
	"Go-Starter-Template/domain"
	"Go-Starter-Template/entities"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	JobRepository interface {
		SearchJob(ctx context.Context, filters domain.JobSearchRequest) ([]entities.Job, error)
		GetJobDetail(ctx context.Context, id string) (entities.Job, error)
		ApplyJob(ctx context.Context, jobApplication entities.JobApplication) error
		GetApplicants(ctx context.Context, jobID uuid.UUID) ([]entities.JobApplication, error)
		CheckCompanyIDFromJob(ctx context.Context, jobID uuid.UUID, userID uuid.UUID) error
		ChangeApplicationStatus(ctx context.Context, jobApplication entities.JobApplication) error
		CheckCompanyIDFromApplication(ctx context.Context, jobApplicationID uuid.UUID, userID uuid.UUID) error
		GetJobApplicationByID(ctx context.Context, jobApplicationID uuid.UUID) (entities.JobApplication, error)
	}
	jobRepository struct {
		db *gorm.DB
	}
)

func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) CheckCompanyIDFromJob(ctx context.Context, jobID uuid.UUID, userID uuid.UUID) error {

	var job entities.Job
	var companyID uuid.UUID

	err := r.db.WithContext(ctx).Where("id = ?", jobID).First(&job).Error

	if err != nil {
		return err
	}

	companyID = job.CompanyID

	var company entities.Companies
	err = r.db.WithContext(ctx).Where("user_id = ?", userID).First(&company).Error

	if err != nil {
		return err
	}

	if companyID != company.ID {
		return domain.ErrCompanyNotFound
	}

	return nil
}

func (r *jobRepository) GetJobDetail(ctx context.Context, id string) (entities.Job, error) {
	var job entities.Job
	err := r.db.WithContext(ctx).Preload("Company.User").Preload("Skills").Where("id = ?", id).First(&job).Error

	if err != nil {
		return entities.Job{}, err
	}

	return job, nil
}

func (r *jobRepository) SearchJob(ctx context.Context, filters domain.JobSearchRequest) ([]entities.Job, error) {

	var jobs []entities.Job
	query := r.db.WithContext(ctx).Preload("Company.User").Preload("Skills").Model(&entities.Job{})

	if filters.Title != "" {
		query = query.Where("title ILIKE ?", "%"+filters.Title+"%")
	}

	if filters.JobType != "" {
		query = query.Where("job_type LIKE ?", "%"+filters.JobType+"%")
	}
	if filters.ExperienceLevel != "" {
		query = query.Where("experience_level LIKE ?", "%"+filters.ExperienceLevel+"%")
	}
	if filters.LocationType != "" {
		query = query.Where("location_type LIKE ?", "%"+filters.LocationType+"%")
	}

	if filters.DatePosted != "" {
		if filters.DatePosted == "Past 24 hours" {
			query = query.Where("created_at >= NOW() - INTERVAL '1 day'")
		} else if filters.DatePosted == "Past Week" {
			query = query.Where("created_at >= NOW() - INTERVAL '1 week'")
		} else if filters.DatePosted == "Past Month" {
			query = query.Where("created_at >= NOW() - INTERVAL '1 month'")
		}
	}

	if filters.SortBy != "" {
		if filters.SortBy == "recent" {
			query = query.Order("created_at DESC")
		} else if filters.SortBy == "salary-high" {
			query = query.Order("salary_max DESC")
		} else if filters.SortBy == "salary-low" {
			query = query.Order("salary_min ASC")
		}
	}

	if filters.MinSalary > 0 && filters.MaxSalary > 0 {
		query = query.Where("salary_min >= ?", filters.MinSalary)
		query = query.Where("salary_min <= ?", filters.MaxSalary)
	}

	err := query.Find(&jobs).Error

	if err != nil {
		return []entities.Job{}, err
	}

	return jobs, nil
}

func (r *jobRepository) ApplyJob(ctx context.Context, jobApplication entities.JobApplication) error {

	if err := r.db.WithContext(ctx).Create(&jobApplication).Error; err != nil {
		return err
	}

	return nil
}

func (r *jobRepository) GetApplicants(ctx context.Context, jobID uuid.UUID) ([]entities.JobApplication, error) {
	var applicants []entities.JobApplication
	err := r.db.WithContext(ctx).Preload("User").Where("job_id = ?", jobID).Find(&applicants).Error

	if err != nil {
		return []entities.JobApplication{}, err
	}

	return applicants, nil
}

func (r *jobRepository) GetJobApplicationByID(ctx context.Context, jobApplicationID uuid.UUID) (entities.JobApplication, error) {
	var jobApplication entities.JobApplication
	err := r.db.WithContext(ctx).Preload("Job.Company").Where("id = ?", jobApplicationID).First(&jobApplication).Error

	if err != nil {
		return entities.JobApplication{}, err
	}

	return jobApplication, nil
}

func (r *jobRepository) ChangeApplicationStatus(ctx context.Context, jobApplication entities.JobApplication) error {
	res := r.db.WithContext(ctx).Model(&entities.JobApplication{}).Where("id = ?", jobApplication.ID).Updates(&jobApplication)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *jobRepository) CheckCompanyIDFromApplication(ctx context.Context, jobApplicationID uuid.UUID, userID uuid.UUID) error {

	var jobApplication entities.JobApplication
	var job entities.Job
	var companyID uuid.UUID

	err := r.db.WithContext(ctx).Where("id = ?", jobApplicationID).First(&jobApplication).Error

	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).Where("id = ?", jobApplication.JobID).First(&job).Error

	if err != nil {
		return err
	}

	companyID = job.CompanyID

	var company entities.Companies
	err = r.db.WithContext(ctx).Where("user_id = ?", userID).First(&company).Error

	if err != nil {
		return err
	}

	if companyID != company.ID {
		return domain.ErrCompanyNotFound
	}

	return nil
}

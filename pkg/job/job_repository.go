package job

import (
	"Go-Starter-Template/domain"
	"Go-Starter-Template/entities"
	"context"

	"gorm.io/gorm"
)

type (
	JobRepository interface {
		SearchJob(ctx context.Context, filters domain.JobSearchRequest) ([]entities.Job, error)
	}
	jobRepository struct {
		db *gorm.DB
	}
)

func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) SearchJob(ctx context.Context, filters domain.JobSearchRequest) ([]entities.Job, error) {

	var jobs []entities.Job
	query := r.db.WithContext(ctx).Preload("Company").Preload("Skills").Model(&entities.Job{})

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

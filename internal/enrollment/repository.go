package enrollment

import (
	"context"
	"fmt"
	"log"

	"github.com/tzincker/gocourse_domain/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, course *domain.Enrollment) (*domain.Enrollment, error)
	GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error)
	Count(ctx context.Context, filters Filters) (int64, error)
}

type repo struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(ctx context.Context, enrollment *domain.Enrollment) (*domain.Enrollment, error) {
	result := repo.db.WithContext(ctx).Create(enrollment)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return nil, result.Error
	}
	repo.log.Println("enrollment created with id: ", enrollment.ID)
	return enrollment, nil
}

func (repo *repo) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error) {
	var enrollment []domain.Enrollment
	tx := repo.db.WithContext(ctx).Model(&enrollment)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at desc").Find(&enrollment)

	if result.Error != nil {
		repo.log.Println(result.Error)
		return nil, result.Error
	}
	repo.log.Println("enrolllment got")
	return enrollment, nil
}

func (repo *repo) Count(ctx context.Context, filters Filters) (int64, error) {
	var count int64
	tx := repo.db.WithContext(ctx).Model(domain.Enrollment{})
	tx = applyFilters(tx, filters)

	if err := tx.Count(&count).Error; err != nil {
		repo.log.Println(err)
		return 0, err
	}

	repo.log.Println("enrolllment count")
	return count, nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.UserId != "" {
		filters.UserId = fmt.Sprintf("%%%s%%", filters.UserId)
		tx = tx.Where("user_id = ?", filters.UserId)
	}

	if filters.CourseId != "" {
		filters.CourseId = fmt.Sprintf("%s", filters.CourseId)
		tx = tx.Where("course_id = ?", filters.CourseId)
	}
	return tx
}

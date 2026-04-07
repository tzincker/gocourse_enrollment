package enrollment

import (
	"context"
	"log"

	"github.com/tzincker/gocourse_domain/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, course *domain.Enrollment) (*domain.Enrollment, error)
	GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error)
	Get(ctx context.Context, id string) (*domain.Enrollment, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, status *string) error
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

func (repo *repo) Get(ctx context.Context, id string) (*domain.Enrollment, error) {
	enrollment := domain.Enrollment{ID: id}

	result := repo.db.WithContext(ctx).First(&enrollment)
	if result.Error != nil {
		repo.log.Println(result.Error)

		if result.Error == gorm.ErrRecordNotFound {
			return nil, &ErrNotFound{EnrollmentId: id}
		}

		return nil, result.Error
	}
	repo.log.Println("enrolllment found with id: ", enrollment.ID)
	return &enrollment, nil
}

func (repo *repo) Delete(ctx context.Context, id string) error {
	enrollment := domain.Enrollment{ID: id}

	result := repo.db.WithContext(ctx).Delete(&enrollment)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		repo.log.Printf("enrollment %s doesn't exists", id)
		return &ErrNotFound{EnrollmentId: id}
	}

	repo.log.Println("enrollment deleted with id: ", enrollment.ID)
	return nil
}

func (repo *repo) Update(
	ctx context.Context,
	id string,
	status *string,
) error {

	values := make(map[string]any)

	if status != nil {
		values["status"] = status
	}

	result := repo.db.WithContext(ctx).Model(&domain.Enrollment{}).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		repo.log.Printf("enrolllment %s doesn't exists", id)
		return &ErrNotFound{EnrollmentId: id}
	}

	repo.log.Println("enrollment updated with id: ", id)
	return nil
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
		tx = tx.Where("user_id = ?", filters.UserId)
	}

	if filters.CourseId != "" {
		tx = tx.Where("course_id = ?", filters.CourseId)
	}
	return tx
}

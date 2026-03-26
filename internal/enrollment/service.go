package enrollment

import (
	"context"
	"log"

	"github.com/tzincker/gocourse_domain/domain"
)

type (
	Service interface {
		Create(ctx context.Context, userId, courseId string) (*domain.Enrollment, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error)
		Count(ctx context.Context, filters Filters) (int64, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(ctx context.Context, userID, courseID string) (*domain.Enrollment, error) {
	log.Println("Create enrollment service")
	enrollment := domain.Enrollment{
		UserID:   userID,
		CourseID: courseID,
		Status:   "P",
	}

	// if _, err := s.userSrv.Get(enrollment.UserID); err != nil {
	// 	return nil, errors.New("user id does not exist")
	// }

	// if _, err := s.courseSrv.Get(enrollment.CourseID); err != nil {
	// 	return nil, errors.New("course id does not exist")
	// }

	e, err := s.repo.Create(ctx, &enrollment)

	if err != nil {
		s.log.Println(err)
	}

	return e, err
}

func (s service) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error) {
	log.Println("Get all enrollments service")

	enrollments, err := s.repo.GetAll(ctx, filters, offset, limit)

	if err != nil {
		s.log.Println(err)
	}

	return enrollments, err
}

func (s service) Count(ctx context.Context, filters Filters) (int64, error) {
	log.Println("Get all enrollments count service")
	count, err := s.repo.Count(ctx, filters)
	if err != nil {
		s.log.Println(err)
	}

	return count, err
}

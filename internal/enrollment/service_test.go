package enrollment_test

import (
	"context"
	"errors"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	courseSdkMock "github.com/tzincker/go_course_sdk/course/mock"
	userSdkMock "github.com/tzincker/go_course_sdk/user/mock"
	"github.com/tzincker/gocourse_domain/domain"
	enrollmentPkg "github.com/tzincker/gocourse_enrollment/internal/enrollment"
)

func TestService_GetAll(t *testing.T) {
	l := log.New(io.Discard, "", 0)

	t.Run("should return an error", func(t *testing.T) {
		var want error = errors.New("mock error")
		var wantCounter = 1
		var count int = 0

		repo := &mockRepository{
			GetAllMock: func(ctx context.Context, filters enrollmentPkg.Filters, offset, limit int) ([]domain.Enrollment, error) {
				count++
				return nil, errors.New("mock error")
			},
		}
		service := enrollmentPkg.NewService(l, nil, nil, repo)

		enrollments, err := service.GetAll(context.Background(), enrollmentPkg.Filters{}, 0, 10)

		assert.Error(t, err)
		assert.Nil(t, enrollments)
		assert.Equal(t, wantCounter, count)
		assert.EqualError(t, want, err.Error())
	})

	t.Run("should return all enrollments", func(t *testing.T) {
		var wantCounter int = 1
		var count int = 0

		want := []domain.Enrollment{
			{
				ID:       "1",
				UserID:   "1",
				CourseID: "1",
				Status:   "A",
			},
		}

		repo := &mockRepository{
			GetAllMock: func(ctx context.Context, filters enrollmentPkg.Filters, offset, limit int) ([]domain.Enrollment, error) {
				count++
				return []domain.Enrollment{
					{
						ID:       "1",
						UserID:   "1",
						CourseID: "1",
						Status:   "A",
					},
				}, nil
			},
		}
		service := enrollmentPkg.NewService(l, nil, nil, repo)

		enrollments, err := service.GetAll(context.Background(), enrollmentPkg.Filters{}, 0, 10)

		assert.NotNil(t, enrollments)
		assert.Nil(t, err)
		assert.Equal(t, wantCounter, count)
		assert.Equal(t, want, enrollments)
	})
}

func TestService_Update(t *testing.T) {
	l := log.New(io.Discard, "", 0)

	t.Run("should return an error", func(t *testing.T) {
		var want error = errors.New("mock error")
		var wantCounter int = 1
		var count int = 0

		repo := &mockRepository{
			UpdateMock: func(ctx context.Context, id string, status *string) error {
				count++
				return errors.New("mock error")
			},
		}

		service := enrollmentPkg.NewService(l, nil, nil, repo)

		status := "A"
		err := service.Update(context.Background(), "1", &status)

		assert.NotNil(t, err)
		assert.Equal(t, wantCounter, count)
		assert.EqualError(t, want, err.Error())
	})

	t.Run("should update an enrollment", func(t *testing.T) {
		var wantCounter int = 1
		var count int = 0
		var wantID string = "1"
		var wantStatus string = "A"

		repo := &mockRepository{
			UpdateMock: func(ctx context.Context, id string, status *string) error {
				count++
				assert.Equal(t, wantID, id)
				assert.NotNil(t, status)
				assert.Equal(t, wantStatus, *status)
				return nil
			},
		}
		service := enrollmentPkg.NewService(l, nil, nil, repo)

		status := "A"
		err := service.Update(context.Background(), "1", &status)
		assert.Nil(t, err)
		assert.Equal(t, wantCounter, count)
	})
}

func TestService_Count(t *testing.T) {
	l := log.New(io.Discard, "", 0)

	t.Run("should return an error", func(t *testing.T) {
		var want error = errors.New("mock error")
		var wantCounter int = 1
		var counter int = 0

		repo := &mockRepository{
			CountMock: func(ctx context.Context, filters enrollmentPkg.Filters) (int64, error) {
				counter++
				return 0, errors.New("mock error")
			},
		}
		service := enrollmentPkg.NewService(l, nil, nil, repo)
		count, err := service.Count(context.Background(), enrollmentPkg.Filters{})

		assert.NotNil(t, err)
		assert.Equal(t, wantCounter, counter)
		assert.EqualError(t, want, err.Error())
		assert.Zero(t, count)
	})

	t.Run("should return count of enrollments", func(t *testing.T) {
		var want int64 = 1
		var wantCounter int = 1
		var counter int = 0

		repo := &mockRepository{
			CountMock: func(ctx context.Context, filters enrollmentPkg.Filters) (int64, error) {
				counter++
				return 1, nil
			},
		}
		service := enrollmentPkg.NewService(l, nil, nil, repo)
		count, err := service.Count(context.Background(), enrollmentPkg.Filters{})

		assert.Nil(t, err)
		assert.Equal(t, wantCounter, counter)
		assert.Equal(t, want, count)
		assert.NotZero(t, count)

	})
}

func TestService_Create(t *testing.T) {
	l := log.New(io.Discard, "", 0)

	t.Run("it should return error in user sdk", func(t *testing.T) {
		var want error = errors.New("mock error")
		var wantCounter int = 1
		var count int = 0

		userSdk := &userSdkMock.UserSdkMock{
			GetMock: func(id string) (*domain.User, error) {
				count++
				return nil, errors.New("mock error")
			},
		}

		service := enrollmentPkg.NewService(l, userSdk, nil, nil)

		enrollment, err := service.Create(context.Background(), "1", "1")

		assert.NotNil(t, err)
		assert.Equal(t, wantCounter, count)
		assert.EqualError(t, want, err.Error())
		assert.Nil(t, enrollment)
	})

	t.Run("it should return error in course sdk", func(t *testing.T) {
		var want error = errors.New("mock error")
		var wantCounter int = 2
		var count int = 0

		userSdk := &userSdkMock.UserSdkMock{
			GetMock: func(id string) (*domain.User, error) {
				count++
				return nil, nil
			},
		}

		courseSdk := &courseSdkMock.CourseSdkMock{
			GetMock: func(id string) (*domain.Course, error) {
				count++
				return nil, errors.New("mock error")
			},
		}

		service := enrollmentPkg.NewService(l, userSdk, courseSdk, nil)

		enrollment, err := service.Create(context.Background(), "1", "1")

		assert.NotNil(t, err)
		assert.Equal(t, wantCounter, count)
		assert.EqualError(t, want, err.Error())
		assert.Nil(t, enrollment)
	})

	t.Run("it should return error in repo", func(t *testing.T) {
		var want error = errors.New("mock error")
		var wantCounter int = 3
		var count int = 0

		userSdk := &userSdkMock.UserSdkMock{
			GetMock: func(id string) (*domain.User, error) {
				count++
				return nil, nil
			},
		}

		courseSdk := &courseSdkMock.CourseSdkMock{
			GetMock: func(id string) (*domain.Course, error) {
				count++
				return nil, nil
			},
		}

		repo := &mockRepository{
			CreateMock: func(ctx context.Context, enroll *domain.Enrollment) (*domain.Enrollment, error) {
				count++
				return nil, errors.New("mock error")
			},
		}

		service := enrollmentPkg.NewService(l, userSdk, courseSdk, repo)

		enrollment, err := service.Create(context.Background(), "1", "1")

		assert.NotNil(t, err)
		assert.Equal(t, wantCounter, count)
		assert.EqualError(t, want, err.Error())
		assert.Nil(t, enrollment)
	})

	t.Run("it should create enrollment", func(t *testing.T) {
		var want *domain.Enrollment = &domain.Enrollment{
			ID:       "1",
			CourseID: "1",
			UserID:   "1",
			Status:   "P",
		}
		var wantCounter int = 3
		var count int = 0

		userSdk := &userSdkMock.UserSdkMock{
			GetMock: func(id string) (*domain.User, error) {
				count++
				return nil, nil
			},
		}

		courseSdk := &courseSdkMock.CourseSdkMock{
			GetMock: func(id string) (*domain.Course, error) {
				count++
				return nil, nil
			},
		}

		repo := &mockRepository{
			CreateMock: func(ctx context.Context, enroll *domain.Enrollment) (*domain.Enrollment, error) {
				count++
				return &domain.Enrollment{
					ID:       "1",
					CourseID: enroll.CourseID,
					UserID:   enroll.UserID,
					Status:   "P",
				}, nil
			},
		}

		service := enrollmentPkg.NewService(l, userSdk, courseSdk, repo)

		enrollment, err := service.Create(context.Background(), "1", "1")

		assert.Nil(t, err)
		assert.Equal(t, wantCounter, count)
		assert.NotNil(t, enrollment)
		assert.Equal(t, want.ID, enrollment.ID)
		assert.Equal(t, want.UserID, enrollment.UserID)
		assert.Equal(t, want.CourseID, enrollment.CourseID)
		assert.Equal(t, want.Status, enrollment.Status)

	})
}

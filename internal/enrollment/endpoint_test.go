package enrollment_test

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	courseSdk "github.com/tzincker/go_course_sdk/course"
	courseSdkMock "github.com/tzincker/go_course_sdk/course/mock"
	userSdk "github.com/tzincker/go_course_sdk/user"
	userSdkMock "github.com/tzincker/go_course_sdk/user/mock"
	"github.com/tzincker/go_lib_response/response"
	"github.com/tzincker/gocourse_domain/domain"
	"github.com/tzincker/gocourse_enrollment/internal/enrollment"
)

func TestCreateEndpoint(t *testing.T) {
	l := log.New(io.Discard, "", 0)

	t.Run("should return bad request when user id is empty", func(t *testing.T) {
		endpoint := enrollment.MakeEndpoints(nil, enrollment.Config{})
		_, err := endpoint.Create(context.Background(), enrollment.CreateReq{})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, enrollment.ErrUserIDRequired, resp.Error())
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})

	t.Run("should return bad request when course id is empty", func(t *testing.T) {
		endpoint := enrollment.MakeEndpoints(nil, enrollment.Config{})
		_, err := endpoint.Create(context.Background(), enrollment.CreateReq{UserID: "123"})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, enrollment.ErrCourseIDRequired, resp.Error())
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})

	obj := []struct {
		tag            string
		repositoryMock enrollment.Repository
		userSdkMock    userSdk.Transport
		courseSdkMock  courseSdk.Transport
		wantErr        error
		wantCode       int
		wantResponse   *domain.Enrollment
	}{
		{
			tag: "should return an error if user returns an unexpected error",
			userSdkMock: &userSdkMock.UserSdkMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, errors.New("unexpected error")
				},
			},
			wantErr:  errors.New("unexpected error"),
			wantCode: http.StatusInternalServerError,
		},
		{
			tag: "should return an error if user does not exist",
			userSdkMock: &userSdkMock.UserSdkMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, userSdk.ErrNotFound{}
				},
			},
			wantErr:  userSdk.ErrNotFound{},
			wantCode: http.StatusNotFound,
		},
		{
			tag: "should return an error if course sdk returns an unexpected error",
			userSdkMock: &userSdkMock.UserSdkMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, nil
				},
			},
			courseSdkMock: &courseSdkMock.CourseSdkMock{
				GetMock: func(id string) (*domain.Course, error) {
					return nil, errors.New("unexpected error")
				},
			},
			wantErr:  errors.New("unexpected error"),
			wantCode: http.StatusInternalServerError,
		},
		{
			tag: "should return an error if course does not exist",
			userSdkMock: &userSdkMock.UserSdkMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, nil
				},
			},
			courseSdkMock: &courseSdkMock.CourseSdkMock{
				GetMock: func(id string) (*domain.Course, error) {
					return nil, courseSdk.ErrNotFound{}
				},
			},
			wantErr:  courseSdk.ErrNotFound{},
			wantCode: http.StatusInternalServerError,
		},
		{
			tag: "should return an error if repository returns an unexpected error",
			userSdkMock: &userSdkMock.UserSdkMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, nil
				},
			},
			courseSdkMock: &courseSdkMock.CourseSdkMock{
				GetMock: func(id string) (*domain.Course, error) {
					return nil, courseSdk.ErrNotFound{}
				},
			},
			repositoryMock: &mockRepository{
				CreateMock: func(ctx context.Context, enrollment *domain.Enrollment) (*domain.Enrollment, error) {
					return nil, errors.New("unexpected error")
				},
			},
			wantErr:  errors.New("unexpected error"),
			wantCode: http.StatusInternalServerError,
		},

		{
			tag: "should return the enrollment",
			userSdkMock: &userSdkMock.UserSdkMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, nil
				},
			},
			courseSdkMock: &courseSdkMock.CourseSdkMock{
				GetMock: func(id string) (*domain.Course, error) {
					return nil, nil
				},
			},
			repositoryMock: &mockRepository{
				CreateMock: func(ctx context.Context, enrollment *domain.Enrollment) (*domain.Enrollment, error) {
					enrollment.ID = "10"
					return enrollment, nil
				},
			},
			wantCode: http.StatusCreated,
			wantResponse: &domain.Enrollment{
				ID:       "10",
				UserID:   "1",
				CourseID: "4",
				Status:   "P",
			},
		},
	}

	for _, obj := range obj {
		t.Run(obj.tag, func(t *testing.T) {
			service := enrollment.NewService(l, obj.userSdkMock, obj.courseSdkMock, obj.repositoryMock)
			endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
			resp, err := endpoint.Create(context.Background(), enrollment.CreateReq{UserID: "1", CourseID: "4"})

			if obj.wantErr != nil {
				assert.NotNil(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NotNil(t, resp)
				assert.Nil(t, err)

				r := resp.(response.Response)
				assert.Equal(t, obj.wantCode, r.StatusCode())
				assert.Empty(t, r.Error())

				enrollment := r.GetData().(*domain.Enrollment)
				assert.Equal(t, obj.wantResponse.ID, enrollment.ID)
				assert.Equal(t, obj.wantResponse.UserID, enrollment.UserID)
				assert.Equal(t, obj.wantResponse.CourseID, enrollment.CourseID)
				assert.Equal(t, obj.wantResponse.Status, enrollment.Status)
			}
		})
	}

}

func TestGetAllEndpoint(t *testing.T) {
	l := log.New(io.Discard, "", 0)

	t.Run("should return an error if count returns an unexpected error", func(t *testing.T) {
		wantErr := errors.New("unexpected error")
		service := enrollment.NewService(l, nil, nil, &mockRepository{
			CountMock: func(ctx context.Context, filter enrollment.Filters) (int64, error) {
				return 0, errors.New("unexpected error")
			},
		})
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
		_, err := endpoint.GetAll(context.Background(), enrollment.GetAllReq{})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, wantErr, err.Error())
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	})

	t.Run("should return an error if meta returns a parsing error", func(t *testing.T) {
		wantErr := errors.New("strconv.Atoi: parsing \"invalid number\": invalid syntax")
		service := enrollment.NewService(l, nil, nil, &mockRepository{
			CountMock: func(ctx context.Context, filter enrollment.Filters) (int64, error) {
				return 3, nil
			},
		})
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{LimPageDef: "invalid number"})
		_, err := endpoint.GetAll(context.Background(), enrollment.GetAllReq{})

		resp := err.(response.Response)
		assert.EqualError(t, wantErr, err.Error())
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	})

	t.Run("should return an error if GetAll repository returns an unexpected error", func(t *testing.T) {
		wantErr := errors.New("unexpected error")
		service := enrollment.NewService(l, nil, nil, &mockRepository{
			CountMock: func(ctx context.Context, filter enrollment.Filters) (int64, error) {
				return 3, nil
			},
			GetAllMock: func(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]domain.Enrollment, error) {
				return nil, errors.New("unexpected error")
			},
		})
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{LimPageDef: "10"})
		_, err := endpoint.GetAll(context.Background(), enrollment.GetAllReq{})

		resp := err.(response.Response)
		assert.EqualError(t, wantErr, err.Error())
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	})

	t.Run("should return the enrollments", func(t *testing.T) {
		wantEnrollments := []domain.Enrollment{
			{ID: "1", UserID: "1", CourseID: "1", Status: "P"},
			{ID: "2", UserID: "2", CourseID: "2", Status: "P"},
			{ID: "3", UserID: "3", CourseID: "3", Status: "P"},
		}
		service := enrollment.NewService(l, nil, nil, &mockRepository{
			CountMock: func(ctx context.Context, filter enrollment.Filters) (int64, error) {
				return 3, nil
			},
			GetAllMock: func(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]domain.Enrollment, error) {
				return []domain.Enrollment{
					{ID: "1", UserID: "1", CourseID: "1", Status: "P"},
					{ID: "2", UserID: "2", CourseID: "2", Status: "P"},
					{ID: "3", UserID: "3", CourseID: "3", Status: "P"},
				}, nil
			},
		})
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{LimPageDef: "10"})
		resp, err := endpoint.GetAll(context.Background(), enrollment.GetAllReq{})
		assert.Nil(t, err)

		r := resp.(response.Response)
		assert.Equal(t, http.StatusOK, r.StatusCode())
		assert.Empty(t, r.Error())

		enrollments := r.GetData().([]domain.Enrollment)
		assert.Equal(t, wantEnrollments, enrollments)
	})
}

func TestUpdateEndpoint(t *testing.T) {
	l := log.New(io.Discard, "", 0)

	t.Run("should return an error if status is empty", func(t *testing.T) {
		endpoint := enrollment.MakeEndpoints(nil, enrollment.Config{})
		status := ""
		_, err := endpoint.Update(context.Background(), enrollment.UpdateReq{Status: &status})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, enrollment.ErrStatusRequired, err.Error())
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})

	t.Run("should return an error if repository returns a not found error", func(t *testing.T) {
		service := enrollment.NewService(l, nil, nil, &mockRepository{
			UpdateMock: func(ctx context.Context, id string, status *string) error {
				return &enrollment.ErrNotFound{EnrollmentId: id}
			},
		})

		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
		status := "A"
		_, err := endpoint.Update(context.Background(), enrollment.UpdateReq{ID: "20", Status: &status})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, &enrollment.ErrNotFound{EnrollmentId: "20"}, resp.Error())
		assert.Equal(t, http.StatusNotFound, resp.StatusCode())
	})

	t.Run("should return an error if repository returns an non expected error", func(t *testing.T) {
		wantErr := errors.New("unexpected error")
		service := enrollment.NewService(l, nil, nil, &mockRepository{
			UpdateMock: func(ctx context.Context, id string, status *string) error {
				return errors.New("unexpected error")
			},
		})

		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
		status := "A"
		_, err := endpoint.Update(context.Background(), enrollment.UpdateReq{ID: "20", Status: &status})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, wantErr, resp.Error())
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	})

	t.Run("should return an error if status is not valid", func(t *testing.T) {
		wantErr := enrollment.ErrStatusNotValid
		service := enrollment.NewService(l, nil, nil, &mockRepository{
			UpdateMock: func(ctx context.Context, id string, status *string) error {
				return enrollment.ErrStatusNotValid
			},
		})

		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
		status := "X"
		_, err := endpoint.Update(context.Background(), enrollment.UpdateReq{ID: "20", Status: &status})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, wantErr, resp.Error())
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})

	t.Run("should return success", func(t *testing.T) {
		service := enrollment.NewService(l, nil, nil, &mockRepository{
			UpdateMock: func(ctx context.Context, id string, status *string) error {
				assert.Equal(t, "20", id)
				assert.NotNil(t, status)
				assert.Equal(t, "A", *status)
				return nil
			},
		})

		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
		status := "A"
		resp, err := endpoint.Update(context.Background(), enrollment.UpdateReq{ID: "20", Status: &status})
		assert.Nil(t, err)

		r := resp.(response.Response)

		assert.Equal(t, http.StatusOK, r.StatusCode())
		assert.Empty(t, r.Error())
		assert.Nil(t, r.GetData())

	})
}

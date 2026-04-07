package enrollment_test

import (
	"context"

	"github.com/tzincker/gocourse_domain/domain"
	enrollmentPkg "github.com/tzincker/gocourse_enrollment/internal/enrollment"
)

type mockRepository struct {
	CreateMock func(ctx context.Context, enroll *domain.Enrollment) (*domain.Enrollment, error)
	GetAllMock func(ctx context.Context, filters enrollmentPkg.Filters, offset, limit int) ([]domain.Enrollment, error)
	GetMock    func(ctx context.Context, id string) (*domain.Enrollment, error)
	UpdateMock func(ctx context.Context, id string, status *string) error
	CountMock  func(ctx context.Context, filters enrollmentPkg.Filters) (int64, error)
}

func (m *mockRepository) Create(ctx context.Context, enroll *domain.Enrollment) (*domain.Enrollment, error) {
	e, err := m.CreateMock(ctx, enroll)
	return e, err
}

func (m *mockRepository) GetAll(ctx context.Context, filters enrollmentPkg.Filters, offset, limit int) ([]domain.Enrollment, error) {
	return m.GetAllMock(ctx, filters, offset, limit)
}

func (m *mockRepository) Get(ctx context.Context, id string) (*domain.Enrollment, error) {
	e, err := m.GetMock(ctx, id)
	return e, err
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	return m.Delete(ctx, id)
}

func (m *mockRepository) Update(
	ctx context.Context,
	id string,
	status *string,
) error {
	return m.UpdateMock(ctx, id, status)
}

func (m *mockRepository) Count(ctx context.Context, filters enrollmentPkg.Filters) (int64, error) {
	c, err := m.CountMock(ctx, filters)
	return c, err
}

package enrollment

import (
	"context"
	"errors"
	"slices"

	courseSdk "github.com/tzincker/go_course_sdk/course"
	userSdk "github.com/tzincker/go_course_sdk/user"
	"github.com/tzincker/go_lib_response/response"
	"github.com/tzincker/gocourse_domain/domain"
	"github.com/tzincker/gocourse_meta/meta"
)

type (
	Controller func(ctx context.Context, request any) (any, error)

	Endpoints struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	GetReq struct {
		ID string
	}

	CreateReq struct {
		UserID   string `json:"user_id"`
		CourseID string `json:"course_id"`
	}

	GetAllReq struct {
		UserID   string
		CourseID string
		Limit    int
		Page     int
	}

	UpdateReq struct {
		ID     string
		Status *string `json:"status"`
	}

	DeleteReq struct {
		ID string
	}

	Config struct {
		LimPageDef string
	}
)

func MakeEndpoints(s Service, config Config) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s, config),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request any) (any, error) {

		req := request.(CreateReq)

		if req.UserID == "" {
			return nil, response.BadRequest(ErrUserIDRequired.Error())
		}

		if req.CourseID == "" {
			return nil, response.BadRequest(ErrCourseIDRequired.Error())
		}

		course, err := s.Create(ctx, req.UserID, req.CourseID)
		if err != nil {
			if _, ok := errors.AsType[userSdk.ErrNotFound](err); ok {
				return nil, response.NotFound(err.Error())
			}

			if _, ok := errors.AsType[courseSdk.ErrNotFound](err); ok {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", course, nil), nil
	}
}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(GetAllReq)
		filters := Filters{
			UserId:   req.UserID,
			CourseId: req.CourseID,
		}

		count, err := s.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, config.LimPageDef)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		courses, err := s.GetAll(ctx, filters, meta.Offset(), meta.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", courses, meta), nil
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(GetReq)

		enrollment, err := s.Get(ctx, req.ID)
		if err != nil {
			if _, ok := errors.AsType[*ErrNotFound](err); ok {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", enrollment, nil), nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request any) (any, error) {

		req := request.(UpdateReq)

		if req.Status != nil && *req.Status == "" {
			return nil, response.BadRequest(ErrStatusRequired.Error())
		}

		validStatuses := []string{
			string(domain.Pending),
			string(domain.Active),
			string(domain.Studying),
			string(domain.Inactive),
		}

		if !slices.Contains(validStatuses, *req.Status) {
			return nil, response.BadRequest(ErrStatusNotValid.Error())
		}

		err := s.Update(ctx, req.ID, req.Status)

		if err != nil {
			if _, ok := errors.AsType[*ErrNotFound](err); ok {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil, nil), nil
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(DeleteReq)
		err := s.Delete(ctx, req.ID)

		if err != nil {
			if _, ok := errors.AsType[*ErrNotFound](err); ok {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil, nil), nil
	}
}

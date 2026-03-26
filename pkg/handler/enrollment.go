package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/tzincker/go_lib_response/response"
	"github.com/tzincker/gocourse_enrollment/internal/enrollment"
)

func NewEnrollmentHTTPServer(ctx context.Context, endpoints enrollment.Endpoints) http.Handler {

	router := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	router.Handle("/enrollments", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateEnrollment,
		encodeResponse,
		opts...,
	)).Methods("POST")

	router.Handle("/enrollments", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllEnrollments,
		encodeResponse,
		opts...,
	)).Methods("GET")

	return router
}

func decodeCreateEnrollment(_ context.Context, r *http.Request) (any, error) {
	var req enrollment.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetAllEnrollments(_ context.Context, r *http.Request) (any, error) {
	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := enrollment.GetAllReq{
		UserID:   v.Get("user_id"),
		CourseID: v.Get("course_id_id"),
		Limit:    limit,
		Page:     page,
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp any) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utd-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utd-8")

	resp, ok := err.(response.Response)

	if !ok {
		newResponse := response.BadRequest("error parsing body")
		w.WriteHeader(newResponse.StatusCode())
		_ = json.NewEncoder(w).Encode(newResponse)
		return
	}

	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)

}

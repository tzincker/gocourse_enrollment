package enrollments_test

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/ncostamagna/go_http_client/client"
	courseSdkMock "github.com/tzincker/go_course_sdk/course/mock"
	userSdkMock "github.com/tzincker/go_course_sdk/user/mock"
	"github.com/tzincker/gocourse_domain/domain"
	"github.com/tzincker/gocourse_enrollment/internal/enrollment"
	"github.com/tzincker/gocourse_enrollment/pkg/bootstrap"
	"github.com/tzincker/gocourse_enrollment/pkg/handler"
)

var cli client.Transport

func TestMain(m *testing.M) {

	_ = godotenv.Load("../.env")
	address := bootstrap.Address()

	log := log.New(io.Discard, "", 0)
	db, err := bootstrap.DBConnection()

	if err != nil {
		log.Fatal(err)
	}

	tx := db.Begin()

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == "" {
		log.Fatal("paginator limit default is required")
	}

	userSdk := &userSdkMock.UserSdkMock{
		GetMock: func(id string) (*domain.User, error) {
			return nil, nil
		},
	}

	courseSdk := &courseSdkMock.CourseSdkMock{
		GetMock: func(id string) (*domain.Course, error) {
			return nil, nil
		},
	}

	ctx := context.Background()
	enrollmentRepo := enrollment.NewRepo(log, tx)
	enrollmentSrv := enrollment.NewService(log, userSdk, courseSdk, enrollmentRepo)

	h := handler.NewEnrollmentHTTPServer(ctx, enrollment.MakeEndpoints(enrollmentSrv, enrollment.Config{LimPageDef: pagLimDef}))

	cli = client.New(http.Header{}, "http://"+address, 60000*time.Millisecond, false)

	srv := &http.Server{
		Handler:      accessControl(h),
		Addr:         address,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	errCh := make(chan error)

	go func() {
		log.Printf("User Server listening to: %s\n", address)
		errCh <- srv.ListenAndServe()
	}()

	r := m.Run()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println(err)
	}

	tx.Rollback()
	os.Exit(r)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS, HEAD, DELETE")

		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

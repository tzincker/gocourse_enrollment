package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/tzincker/gocourse_enrollment/internal/enrollment"
	"github.com/tzincker/gocourse_enrollment/pkg/bootstrap"
	"github.com/tzincker/gocourse_enrollment/pkg/handler"
)

func main() {
	_ = godotenv.Load()
	address := bootstrap.Address()

	log := bootstrap.InitLogger()
	db, err := bootstrap.DBConnection()

	if err != nil {
		log.Fatal(err)
	}

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == "" {
		log.Fatal("paginator limit default is required")
	}

	ctx := context.Background()
	enrollmentRepo := enrollment.NewRepo(log, db)
	enrollmentSrv := enrollment.NewService(log, enrollmentRepo)

	h := handler.NewEnrollmentHTTPServer(ctx, enrollment.MakeEndpoints(enrollmentSrv, enrollment.Config{LimPageDef: pagLimDef}))

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

	err = <-errCh

	if err != nil {
		log.Fatal(err)
	}

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

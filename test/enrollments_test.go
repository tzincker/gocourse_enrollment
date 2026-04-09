package enrollments_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tzincker/gocourse_domain/domain"
	"github.com/tzincker/gocourse_enrollment/internal/enrollment"
)

type dataResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    any    `json:"data"`
	Meta    any    `json:"meta"`
}

func TestEnrollments(t *testing.T) {

	t.Run("should create an enrollment and get it", func(t *testing.T) {
		bodyRequest := &enrollment.CreateReq{
			UserID:   "11-test",
			CourseID: "12-test",
		}

		resp := cli.Post("/enrollments", bodyRequest)
		assert.Nil(t, resp.Err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		dataCreated := domain.Enrollment{}
		dRespCreated := dataResponse{Data: &dataCreated}
		err := resp.FillUp(&dRespCreated)
		assert.Nil(t, err)

		assert.Equal(t, "success", dRespCreated.Message)
		assert.NotEmpty(t, dataCreated.ID)
		assert.Equal(t, http.StatusCreated, dRespCreated.Status)
		assert.Equal(t, bodyRequest.UserID, dataCreated.UserID)
		assert.Equal(t, bodyRequest.CourseID, dataCreated.CourseID)
		assert.Equal(t, "P", dataCreated.Status)

		resp = cli.Get("/enrollments?user_id=" + dataCreated.UserID + "&course_id=" + dataCreated.CourseID)
		assert.Nil(t, resp.Err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var dataGetAll []domain.Enrollment
		dRespGetAll := dataResponse{Data: &dataGetAll}
		err = resp.FillUp(&dRespGetAll)
		assert.Equal(t, 1, len(dataGetAll))
		assert.Equal(t, dataCreated.ID, dataGetAll[0].ID)
		assert.Equal(t, dataCreated.UserID, dataGetAll[0].UserID)
		assert.Equal(t, dataCreated.CourseID, dataGetAll[0].CourseID)
		assert.Equal(t, dataCreated.Status, dataGetAll[0].Status)

		assert.Nil(t, err)

	})

	t.Run("update an enrollment", func(t *testing.T) {
		bodyRequest := &enrollment.CreateReq{
			UserID:   "21-test",
			CourseID: "22-test",
		}

		resp := cli.Post("/enrollments", bodyRequest)
		assert.Nil(t, resp.Err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		dataCreated := domain.Enrollment{}
		dRespCreated := dataResponse{Data: &dataCreated}
		err := resp.FillUp(&dRespCreated)
		assert.Nil(t, err)

		resp = cli.Get("/enrollments?user_id=" + dataCreated.UserID + "&course_id=" + dataCreated.CourseID)
		assert.Nil(t, resp.Err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var dataGetAll []domain.Enrollment
		dRespGetAll := dataResponse{Data: &dataGetAll}
		err = resp.FillUp(&dRespGetAll)
		assert.Equal(t, 1, len(dataGetAll))
		assert.Equal(t, dataCreated.ID, dataGetAll[0].ID)
		assert.Equal(t, dataCreated.UserID, dataGetAll[0].UserID)
		assert.Equal(t, dataCreated.CourseID, dataGetAll[0].CourseID)
		assert.Equal(t, "P", dataGetAll[0].Status)
		assert.Nil(t, err)

		status := "A"
		updateRequest := enrollment.UpdateReq{
			Status: &status,
		}
		resp = cli.Patch("/enrollments/"+dataCreated.ID, updateRequest)
		assert.Nil(t, resp.Err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		dRespUpdated := dataResponse{Data: &dataCreated}
		err = resp.FillUp(&dRespUpdated)
		assert.Nil(t, err)

		resp = cli.Get("/enrollments?user_id=" + dataCreated.UserID + "&course_id=" + dataCreated.CourseID)
		assert.Nil(t, resp.Err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		dRespGetAll = dataResponse{Data: &dataGetAll}
		err = resp.FillUp(&dRespGetAll)
		assert.Equal(t, 1, len(dataGetAll))
		assert.Equal(t, dataCreated.ID, dataGetAll[0].ID)
		assert.Equal(t, dataCreated.UserID, dataGetAll[0].UserID)
		assert.Equal(t, dataCreated.CourseID, dataGetAll[0].CourseID)
		assert.Equal(t, "A", dataGetAll[0].Status)
		assert.Nil(t, err)
	})
}

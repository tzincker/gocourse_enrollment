package enrollment

import (
	"errors"
	"fmt"
)

var ErrUserIDRequired = errors.New("user_id is required")
var ErrCourseIDRequired = errors.New("course_id is required")
var ErrStatusRequired = errors.New("status is required")
var ErrStatusNotValid = errors.New("status is not valid")

type ErrNotFound struct {
	EnrollmentId string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("enrolllment with id: '%s' doesn't exist", e.EnrollmentId)
}

type ErrCourseNotFound struct {
	CourseId string
}

func (e *ErrCourseNotFound) Error() string {
	return fmt.Sprintf("course '%s' doesn't exist", e.CourseId)
}

type ErrUserNotFound struct {
	UserId string
}

func (e *ErrUserNotFound) Error() string {
	return fmt.Sprintf("user '%s' doesn't exist", e.UserId)
}

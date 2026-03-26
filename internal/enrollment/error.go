package enrollment

import (
	"errors"
	"fmt"
)

var ErrUserIDRequired = errors.New("user_id is required")
var ErrCourseIDRequired = errors.New("course_id is required")

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

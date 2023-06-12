package sql_errors

import (
	"strings"

	"github.com/pkg/errors"
)

type SqlNotFoundError struct {
	err error
}

var SqlNotFound *SqlNotFoundError

func (s *SqlNotFoundError) Error() string {
	return errors.Wrap(s.err, "row not found").Error()
}

func NewSqlNotFoundError(err error) error {
	return &SqlNotFoundError{
		err: err,
	}
}

type SqlError struct{}

func NewSqlError() error {
	return &SqlError{}
}

func (s *SqlError) Error() string {
	return "sql error"
}

func ParseSqlErrors(err error) error {
	if strings.Contains(err.Error(), "no rows in result set") {
		return NewSqlNotFoundError(err)
	}
	return NewSqlError()
}

package commonerrors

import (
	"errors"
	"fmt"
)

var ErrSQLRequestNoExist = errors.New("sql request no exist")

type ErrorSQLRequestNoExist struct {
	filename string
}

func (e *ErrorSQLRequestNoExist) Error() string {
	return fmt.Sprintf("user \"%v\" no exist", e.filename)
}

func (e *ErrorSQLRequestNoExist) Unwrap() error {
	return ErrSQLRequestNoExist
}

func NewErrorErrorSQLRequestNoExist(filename string) error {
	return &ErrorSQLRequestNoExist{
		filename: filename,
	}
}

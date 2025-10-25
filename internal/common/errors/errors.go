package commonerrors

import (
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/s-turchinskiy/EffectiveMobile/internal/middleware/logger"
	"go.uber.org/zap"
	"net/http"
	"runtime"
)

const (
	TextErrorGettingData        = "error getting data"
	ContentTypeTextPlainCharset = "text/plain; charset=utf-8"
)

var ErrDuplicateKey = fmt.Errorf("record with such key already exists")

func WrapError(err error) error {

	if err == nil {
		return nil
	}

	_, filename, line, _ := runtime.Caller(1)
	return fmt.Errorf("[error] %s %d: %w", filename, line, err)
}

func ErrorGettingData(w http.ResponseWriter, err error) {

	if err == nil {
		logger.Log.Fatalw("error is empty")
	}

	logger.Log.Infow(TextErrorGettingData, zap.Error(err))
	w.Header().Set("Content-Type", ContentTypeTextPlainCharset)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}

func CheckResponseStatus(statusCode int, body []byte, url string) error {

	if statusCode != http.StatusOK {

		logger.Log.Infow("error. status code <> 200",
			"status code", statusCode,
			"url", url,
			"body", string(body))
		err := fmt.Errorf("status code <> 200, = %d, url : %s", statusCode, url)
		return err
	}

	return nil
}

func IsConnectionError(err error) bool {

	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	errors.As(err, &pgErr)
	if pgErr == nil {
		return false
	}

	return pgerrcode.IsConnectionException(pgErr.Code)

}

func IsDuplicateKeyError(err error) bool {

	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	errors.As(err, &pgErr)
	if pgErr == nil {
		return false
	}

	return pgErr.Code == "23505"

}

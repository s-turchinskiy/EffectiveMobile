package handlers

import (
	"context"
	"github.com/s-turchinskiy/EffectiveMobile/internal/middleware/logger"
	"github.com/s-turchinskiy/EffectiveMobile/internal/repository"
	"github.com/s-turchinskiy/EffectiveMobile/internal/service"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	ContentTypeTextHTML         = "text/html; charset=utf-8"
	ContentTypeTextPlain        = "text/plain"
	ContentTypeTextPlainCharset = "text/plain; charset=utf-8"
	ContentTypeApplicationJSON  = "application/json"
)

type Handler struct {
	Service service.Servicer
}

func NewHandler(ctx context.Context, rep repository.Repository, retryStrategy []time.Duration) *Handler {

	return &Handler{
		Service: service.New(rep, retryStrategy),
	}
}

func internalError(w http.ResponseWriter, err error) {
	logger.Log.Warnw("internal error", zap.Error(err))
	w.Header().Set("Content-Type", ContentTypeTextPlainCharset)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

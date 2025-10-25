package handlers

import (
	"github.com/mailru/easyjson"
	"github.com/s-turchinskiy/EffectiveMobile/internal/middleware/logger"
	"go.uber.org/zap"
	"net/http"
)

// ReadSubscriptions godoc
// @Tags Subscription
// @Summary Чтение подписок
// @Description Чтение всех существующих подписок
// @ID subscriptionReadSubscriptions
// @Accept  json
// @Produce json
// @Success 200 {object} models.Subscriptions "OK"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Внутренняя ошибка"
// @Router /subscription/read [get]
func (h *Handler) ReadSubscriptions(w http.ResponseWriter, r *http.Request) {

	context := r.Context()
	data, err := h.Service.GetSubscriptions(context)
	if err != nil {
		internalError(w, err)
		return
	}

	if len(data) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	rawBytes, err := easyjson.Marshal(data)
	if err != nil {
		logger.Log.Info("error encoding response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", ContentTypeApplicationJSON)
	w.Write(rawBytes)

}

package handlers

import (
	"github.com/s-turchinskiy/EffectiveMobile/internal/middleware/logger"
	"github.com/s-turchinskiy/EffectiveMobile/internal/middleware/read_validate_json"
	"github.com/s-turchinskiy/EffectiveMobile/internal/models"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

// SumSubscriptions godoc
// @Tags Subscriptions
// @Summary Сумма подписок
// @Description Суммарная стоимость всех подписок за выбранный период с фильтрацией по id пользователя и названию подписки
// @ID subscriptionsSumSubscriptions
// @Accept  json
// @Produce html
// @Param data body models.SumSubscriptionsJSON true "Фильтры"
// @Success 200 {object} uint64 "OK"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Внутренняя ошибка"
// @Router /subscription/sum [POST]
func (h *Handler) SumSubscriptions(w http.ResponseWriter, r *http.Request) {

	req := r.Context().Value(readvalidatejson.UnmarshalJSON).(models.SumSubscriptionsJSON)

	startDate := time.Time(*req.Period)

	data := models.SumSubscriptions{
		Period:      &startDate,
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
	}

	sum, err := h.Service.SumSubscriptions(r.Context(), data)
	if err != nil {
		logger.Log.Info(zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(strconv.FormatUint(sum, 10)))

}

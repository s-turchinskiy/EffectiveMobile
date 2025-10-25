package handlers

import (
	"fmt"
	"github.com/s-turchinskiy/EffectiveMobile/internal/middleware/logger"
	readvalidatejson "github.com/s-turchinskiy/EffectiveMobile/internal/middleware/read_validate_json"
	"github.com/s-turchinskiy/EffectiveMobile/internal/models"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

// UpdateSubscription godoc
// @Tags Subscription
// @Summary Обновление подписки
// @Description Обновление существующей подписки на основе данных json
// @ID subscriptionUpdateSubscription
// @Accept  json
// @Produce json
// @Param data body models.UpdateSubscriptionJSON true "Содержимое подписки"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Внутренняя ошибка"
// @Router /subscription/update [PUT]
func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {

	req := r.Context().Value(readvalidatejson.UnmarshalJSON).(models.UpdateSubscriptionJSON)

	id, err := strconv.ParseUint(req.ID, 10, 64)
	if err != nil {

		err = fmt.Errorf("id is not positive integer")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	startDate := time.Time(*req.StartDate)

	data := models.UpdateSubscription{
		ID:          id,
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   &startDate,
	}

	if req.EndDate != nil {
		endDate := time.Time(*req.EndDate)
		data.EndDate = &endDate
	}

	err = h.Service.UpdateSubscription(r.Context(), data)
	if err != nil {
		logger.Log.Info(zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

}

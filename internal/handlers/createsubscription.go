package handlers

import (
	"github.com/s-turchinskiy/EffectiveMobile/internal/middleware/logger"
	readvalidatejson "github.com/s-turchinskiy/EffectiveMobile/internal/middleware/read_validate_json"
	"github.com/s-turchinskiy/EffectiveMobile/internal/models"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// CreateSubscription godoc
// @Tags Subscription
// @Summary Создание новой подписки
// @Description Создание новой подписки на основе данных json
// @ID subscriptionCreateSubscription
// @Accept  json
// @Produce json
// @Param data body models.CreateSubscriptionJSON true "Содержимое подписки"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Внутренняя ошибка"
// @Router /subscription/create [POST]
func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {

	req := r.Context().Value(readvalidatejson.UnmarshalJSON).(models.CreateSubscriptionJSON)
	startDate := time.Time(*req.StartDate)

	data := models.CreateSubscription{
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   &startDate,
	}

	if req.EndDate != nil {
		endDate := time.Time(*req.EndDate)
		data.EndDate = &endDate
	}

	err := h.Service.CreateSubscription(r.Context(), data)
	if err != nil {
		logger.Log.Info(zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

}

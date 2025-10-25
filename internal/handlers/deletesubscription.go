package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

// DeleteSubscription godoc
// @Tags Subscription
// @Summary Удаление подписки
// @Description Удаление подписки по ID
// @ID subscriptionDeleteSubscription
// @Accept  json
// @Produce html
// @Param Id path uint64 true "Id"
// @Success 200 {string} string "ОК"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Внутренняя ошибка"
// @Router /subscription/delete{Id} [DELETE]
func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	idStr := query.Get("id")
	if idStr == "" {

		err := fmt.Errorf("id is not defined")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {

		err = fmt.Errorf("id is not positive integer")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	err = h.Service.DeleteSubscription(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
}

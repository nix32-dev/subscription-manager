package handler

import (
	"net/http"
	"subscription/internal/repository"
	"subscription/internal/service"
)

// @Summary Delete subscription
// @Description DeleteSubscription delete a subscription (NEED ID or UUID!)
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id query string false "Subscription ID" example(123)
// @Param uuid query string false "Subscription UUID" format(uuid) example(60601fee-2bf1-4721-ae6f-7636e79a0cba)
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /delete [delete]
func DeleteSubscriptionH(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(service.WrongMethod.Error()))
		service.Logger.Error(service.WrongMethod.Error())
		return
	}

	var getType string
	if inputUUID := r.URL.Query().Get("uuid"); inputUUID != "" {
		getType = "user_id"
		if err := repository.DeleteSubscription(getType, inputUUID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			service.Logger.Error(err.Error())
			return
		}
		w.Write([]byte("Done!"))
		return
	}
	if inputID := r.URL.Query().Get("id"); inputID != "" {
		getType = "id"
		if err := repository.DeleteSubscription(getType, inputID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			service.Logger.Error(err.Error())
			return
		}
		w.Write([]byte("Done!"))
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(service.NotExistsID_UUID.Error()))
		service.Logger.Error(service.NotExistsID_UUID.Error())
		return
	}
}

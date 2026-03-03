package handler

import (
	"encoding/json"
	"net/http"
	"subscription/internal/model"
	"subscription/internal/repository"
	"subscription/internal/service"
)

// @Summary Change subscription
// @Description ChangeSubscription changes a subscription (NEED ID or UUID!)
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id query string false "Subscription ID" example(123)
// @Param uuid query string false "Subscription UUID" format(uuid) example(60601fee-2bf1-4721-ae6f-7636e79a0cba)
// @Param subscription body model.Subscription true "Optional: service_name, price, start_date, end_date"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /change [patch]
func ChangeSubscriptionH(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(service.WrongMethod.Error()))
		service.Logger.Error(service.WrongMethod.Error())
		return
	}

	var input model.Subscription
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		service.Logger.Error(err.Error())
		return
	}

	var getType string
	if inputUUID := r.URL.Query().Get("uuid"); inputUUID != "" {
		getType = "user_id"
		db_changes, err := input.ValidateInputExists(getType, &repository.DataBaseConn, inputUUID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			service.Logger.Error(err.Error())
			return
		}
		if err := repository.ChangeSubscription(db_changes, getType, inputUUID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			service.Logger.Error(err.Error())
			return
		}
		w.Write([]byte("Done!"))

	}
	if inputID := r.URL.Query().Get("id"); inputID != "" {
		getType = "id"
		db_changes, err := input.ValidateInputExists(getType, &repository.DataBaseConn, inputID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			service.Logger.Error(err.Error())
			return
		}
		if err := repository.ChangeSubscription(db_changes, getType, inputID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			service.Logger.Error(err.Error())
			return
		}
		w.Write([]byte("Done!"))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(service.NotExistsID_UUID.Error()))
		service.Logger.Error(service.NotExistsID_UUID.Error())
		return
	}
}

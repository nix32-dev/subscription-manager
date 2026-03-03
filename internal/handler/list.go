package handler

import (
	"encoding/json"
	"net/http"
	"subscription/internal/repository"
	"subscription/internal/service"
)

// @Summary Subscriptions list
// @Description Show subscription by page number (10 subs on one page), uuid/id
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param page query string false "Page number" example(1)
// @Param id query string false "Subscription ID" example(123)
// @Param uuid query string false "Subscription UUID" format(uuid) example(60601fee-2bf1-4721-ae6f-7636e79a0cba)
// @Success 200 {object} model.Subscription
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /list [get]
func GetSubscriptionH(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(service.WrongMethod.Error()))
		service.Logger.Error(service.WrongMethod.Error())
		return
	}

	if r.URL.Query().Get("page") == "" && r.URL.Query().Get("uuid") == "" && r.URL.Query().Get("id") == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(service.WrongPageNumber.Error()))
		service.Logger.Error(service.WrongPageNumber.Error())
		return
	}

	if page := r.URL.Query().Get("page"); page != "" {
		output, err := repository.GetSubscriptions(page)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			service.Logger.Error(err.Error())
			return
		}
		if outputW, err := json.Marshal(output); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			service.Logger.Error(err.Error())
			return
		} else {
			service.Logger.Info("Subscriptions sent!")
			w.Write(outputW)
			return
		}
	}

	var getType string
	var getID string
	if inputUUID := r.URL.Query().Get("uuid"); inputUUID != "" {
		getType = "user_id"
		getID = inputUUID
	} else if inputID := r.URL.Query().Get("id"); inputID != "" {
		getType = "id"
		getID = inputID
	}

	output, err := repository.GetSubscription(getID, getType)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		service.Logger.Error(err.Error())
		return
	}
	if outputW, err := json.Marshal(output); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		service.Logger.Error(err.Error())
		return
	} else {
		service.Logger.Info("Subscription details sent!")
		w.Write(outputW)
		return
	}
}

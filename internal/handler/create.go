package handler

import (
	"encoding/json"
	"net/http"
	"subscription/internal/model"
	"subscription/internal/repository"
	"subscription/internal/service"
)

// @Summary Create subscription
// @Description CreateSubscription creates a new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body model.Subscription true "Need: service_name, price, start_date, end_date(optional)"
// @Success 201 {object} model.Subscription
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /create [post]
func CreateSubscriptionH(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	if err, output := repository.CreateSubscription(input); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		service.Logger.Error(err.Error())
		return
	} else {
		outputW, err := json.Marshal(&output)
		if err != nil {
			service.Logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		service.Logger.Info("Subscription created successfully! UUID: " + output.User_id)
		w.WriteHeader(http.StatusCreated)
		w.Write(outputW)
	}
}

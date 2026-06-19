package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	//	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5"
	"github.com/vivaldi7/golang_code/subscription_service/innternal/models"
	"github.com/vivaldi7/golang_code/subscription_service/innternal/service"
	"github.com/vivaldi7/golang_code/subscription_service/logger"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
}

func NewSubscriptionHandler(service *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	subscription, err := h.service.Create(r.Context(), &req)
	if err != nil {
		logger.Log.Error("Failed to create subscription", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subscription)
}

func (h *SubscriptionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	iD, _ := strconv.Atoi(id)

	subscription, err := h.service.Get(r.Context(), iD)
	if err != nil {
		logger.Log.Error("Failed to get subscription", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if subscription == nil {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscription)

}

func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page-size"))

	subscription, err := h.service.List(r.Context(), int64(page), int64(pageSize))
	if err != nil {
		logger.Log.Error("Failed to list subscriptions", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscription)

}

func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	iD, _ := strconv.Atoi(id)

	var req models.UpdateSubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.Update(r.Context(), iD, &req)
	if err != nil {
		logger.Log.Error("Failed to update subscription", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	iD, _ := strconv.Atoi(id)

	err := h.service.Delete(r.Context(), iD)
	if err != nil {
		logger.Log.Error("Failed to delete subscription", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) GetTotalCost(w http.ResponseWriter, r *http.Request) {
	var req models.TotalCostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	total, err := h.service.GetTotalCost(r.Context(), &req)
	if err != nil {
		logger.Log.Error("Failed to calculate total cost", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]int{"total_cost": total}
	w.Header().Set("Content_Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

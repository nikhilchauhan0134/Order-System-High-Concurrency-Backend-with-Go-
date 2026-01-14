package handlers

import (
	"OrderSystemHighConcurrency/order-api/internal/contracts"
	"OrderSystemHighConcurrency/shared/models"
	"encoding/json"
	"net/http"
)

// OrderHandler handles HTTP requests for orders
type OrderHandler struct {
	orderService contracts.OrderService
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(service contracts.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: service,
	}
}

// ServeHTTP handles POST /orders
func (h *OrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var order models.Order

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields (optional, basic example)
	if order.OrderID == "" || order.Amount <= 0 {
		http.Error(w, "invalid order data", http.StatusBadRequest)
		return
	}

	// Call the service to create the order
	if err := h.orderService.CreateOrder(r.Context(), &order); err != nil {
		http.Error(w, "failed to create order", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusAccepted)
	resp := map[string]string{
		"status":  "order accepted",
		"orderId": order.OrderID,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

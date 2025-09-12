package http

import (
	"encoding/json"
	"net/http"
	"rest-api-in-gin/internal/payment/application"
	"rest-api-in-gin/internal/payment/application/dto"
	"rest-api-in-gin/internal/payment/application/service"
)

// PaymentHandler handles all request for Payment
type PaymentHandler struct {
	depositService *service.DepositFundService
}

// NewPaymentHandler creates a new instance of PaymentHandler
func NewPaymentHandler(service *service.DepositFundService) *PaymentHandler {
	return &PaymentHandler{service}
}

// Deposit handles POST requests to deposit an amount into a Wallet
func (handler *PaymentHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.DepositFundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	userDTO, err := handler.depositService.DepositFund(&req)
	if err != nil {
		mapAppErrorToHTTP(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userDTO)
}

// mapAppErrorToHTTP maps application-layer errors to appropriate HTTP response code
func mapAppErrorToHTTP(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*application.AppError); ok {
		switch appErr.Code() {
		case application.ErrInvalidInput:
			http.Error(w, appErr.Error(), http.StatusBadRequest)
		case application.ErrInternal:
			http.Error(w, appErr.Error(), http.StatusInternalServerError)
		default:
			http.Error(w, "unknown application error", http.StatusServiceUnavailable)
		}
		return
	}

	http.Error(w, "internal server error", http.StatusInternalServerError)
}

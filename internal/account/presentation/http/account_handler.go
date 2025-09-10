package http

import (
	"encoding/json"
	"net/http"

	"rest-api-in-gin/internal/account/application"
	"rest-api-in-gin/internal/account/application/dto"
	"rest-api-in-gin/internal/account/application/service"
)

// AccountHandler handles all request for UserAccount
type AccountHandler struct {
	registerService *service.RegisterAccountService
}

// NewAccountHandler creates a new instance of AccountHandler
func NewAccountHandler(service *service.RegisterAccountService) *AccountHandler {
	return &AccountHandler{service}
}

// Register handles POST requests to register a new UserAccount
func (handler *AccountHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.RegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	userDTO, err := handler.registerService.Register(&req)
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

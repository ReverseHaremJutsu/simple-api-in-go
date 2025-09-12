package http

import "net/http"

// Route defines a simple contract for a single endpoint
type Route struct {
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

// PaymentRoutes exposes routes owned by the payment module
func PaymentRoutes(handler *PaymentHandler) []Route {
	return []Route{
		{Method: http.MethodPost, Pattern: "/payments/deposit", Handler: handler.Deposit},
		// Later:
		// {Method: http.MethodPut, Pattern: "/accounts/{id}", Handler: h.Update},
		// {Method: http.MethodDelete, Pattern: "/accounts/{id}", Handler: h.Delete},
	}
}

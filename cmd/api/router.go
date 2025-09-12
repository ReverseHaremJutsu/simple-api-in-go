package main

import (
	"net/http"

	accounthttp "rest-api-in-gin/internal/account/presentation/http"
	// order "rest-api-in-gin/internal/order/presentation/http"
	paymenthttp "rest-api-in-gin/internal/payment/presentation/http"
	// event "rest-api-in-gin/internal/event/presentation/http"
)

const apiVersion = "/api/v1"

// SetupRouter wires all the routes from all modules
func SetupRouter(accountHandler *accounthttp.AccountHandler, paymentHandler *paymenthttp.PaymentHandler) http.Handler {
	mux := http.NewServeMux()

	for _, r := range accounthttp.AccountRoutes(accountHandler) {
		mux.Handle(apiVersion+r.Pattern, r.Handler)
	}

	for _, r := range paymenthttp.PaymentRoutes(paymentHandler) {
		mux.Handle(apiVersion+r.Pattern, r.Handler)
	}

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	return mux
}

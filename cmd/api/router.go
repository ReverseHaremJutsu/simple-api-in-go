package main

import (
	"net/http"

	accounthttp "rest-api-in-gin/internal/account/presentation/http"
	// order "rest-api-in-gin/internal/order/presentation/http"
	// payment "rest-api-in-gin/internal/payment/presentation/http"
	// event "rest-api-in-gin/internal/event/presentation/http"
)

const apiVersion = "/api/v1"

func SetupRouter(accountHandler *accounthttp.AccountHandler) http.Handler {
	mux := http.NewServeMux()

	// Register bounded context routes
	for _, r := range accounthttp.AccountRoutes(accountHandler) {
		mux.Handle(apiVersion+r.Pattern, r.Handler)
	}

	// health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	return mux
}

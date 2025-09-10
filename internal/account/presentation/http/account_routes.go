package http

import "net/http"

// Route defines a simple contract for a single endpoint
type Route struct {
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

// AccountRoutes exposes routes owned by the account module
func AccountRoutes(handler *AccountHandler) []Route {
	return []Route{
		{Method: http.MethodPost, Pattern: "/accounts", Handler: handler.Register},
		// Later:
		// {Method: http.MethodPut, Pattern: "/accounts/{id}", Handler: h.Update},
		// {Method: http.MethodDelete, Pattern: "/accounts/{id}", Handler: h.Delete},
	}
}

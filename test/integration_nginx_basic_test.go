package test

// Distinguish your calls to fatalf and errorf

import (
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	mock_backend "rest-api-in-gin/internal/mock"
)

const nginxURL = "http://localhost"

func startMockBackend() {
	r := mock_backend.NewMockServer()
	// Run in a goroutine so it doesn't block the test
	go func() {
		if err := r.Run(":8080"); err != nil {
			panic(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
}

func TestNginxReverseProxy(t *testing.T) {

	// Arrange
	startMockBackend()

	tests := map[string]struct {
		path       string
		wantStatus int
		check      func(t *testing.T, resp *http.Response, body string)
	}{
		"Reverse proxy is returning responses": {
			path:       "/api/events/1",
			wantStatus: http.StatusOK,
			check: func(t *testing.T, resp *http.Response, body string) {
				if got := resp.Header.Get("X-Nginx-Server"); got != "my-nginx" {
					t.Errorf("Expected X-Nginx-Server=my-nginx, got %q", got)
				}
			},
		},
		"X-Real-IP is forwarded": {
			path:       "/api/echo-realip",
			wantStatus: http.StatusOK,
			check: func(t *testing.T, resp *http.Response, body string) {
				if body == "" {
					t.Errorf("Expected X-Real-IP header to be forwarded, got empty")
				}

				if net.ParseIP(body) == nil {
					t.Errorf("Expected valid IP address in X-Real-IP, got %q", body)
				}

			},
		},
		"X-Forwarded-For is forwarded": {
			path:       "/api/echo-xff",
			wantStatus: http.StatusOK,
			check: func(t *testing.T, resp *http.Response, body string) {
				if body == "" {
					t.Errorf("Expected X-Forwarded-For header to be forwarded, got empty")
				}

				ips := strings.Split(body, ",")
				if len(ips) == 0 || net.ParseIP(strings.TrimSpace(ips[0])) == nil {
					t.Errorf("Expected valid IP in X-Forwarded-For, got %q", body)
				}
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			// Act
			// Get response from mock backend server
			url := nginxURL + tt.path
			resp, err := http.Get(url)
			if err != nil {
				t.Fatalf("Failed to GET %s: %v", url, err)
			}
			// resp.Body (inside http.Response) is an io.ReadCloser which is a stream connected to the response
			defer resp.Body.Close()

			// Read response body received
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to receive response: %v", err)
			}
			body := string(bodyBytes)

			// Assert
			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("Expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			if tt.check != nil {
				tt.check(t, resp, body)
			}
		})
	}
}

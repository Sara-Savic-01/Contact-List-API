package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name         string
		authHeader   string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Valid token",
			authHeader:   "Bearer Axf2FVAusahoXmKMLZih7LrhBwmYLVmyLDiMoYizPGReJTKEaseAb12oGYvbLleS",
			expectedCode: http.StatusOK,
			expectedBody: "OK",
		},
		{
			name:         "Invalid token",
			authHeader:   "Bearer invalid token",
			expectedCode: http.StatusForbidden,
			expectedBody: "Forbidden\n",
		},
		{
			name:         "Missing Authorization header",
			authHeader:   "",
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Unauthorized\n",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})
			token := "Axf2FVAusahoXmKMLZih7LrhBwmYLVmyLDiMoYizPGReJTKEaseAb12oGYvbLleS"
			middleware := AuthMiddleware(token, handler)
			req, err := http.NewRequest("GET", "/lists", nil)
			if err != nil {
				t.Fatalf("Could not create request:%v", err)
			}
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()
			middleware.ServeHTTP(rr, req)
			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, status)
			}
			if body := rr.Body.String(); body != tt.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tt.expectedBody, body)
			}
		})
	}

}

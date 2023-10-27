package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/routes"

	mockserver "github.com/anthonydip/flutter-messenger-go/app/storefront-api/webserver/mock"
	mockauth "github.com/anthonydip/flutter-messenger-go/pkg/authentication/mock"

	"github.com/gorilla/mux"
)

func TestAuthentication(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		authResult     mockauth.Result
		expectedStatus int
	}{
		"Auth Passes": {
			expectedStatus: http.StatusOK,
		},
		"Auth Fails": {
			authResult:     mockauth.ValidateJWTFail(),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			srv := mockserver.New().WithAuthentication(test.authResult)

			req, err := http.NewRequest(http.MethodGet, "/ping", nil)
			if err != nil {
				t.Fatalf("test failed while creating new HTTP request, %s", err.Error())
			}

			r := mux.NewRouter()
			r.Use(Authentication(srv))

			r.HandleFunc("/ping", routes.Ping(srv)).Methods(http.MethodGet)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != test.expectedStatus {
				t.Errorf("Expected status code `%d` but got `%03d`.", test.expectedStatus, rr.Code)
			}
		})
	}
}

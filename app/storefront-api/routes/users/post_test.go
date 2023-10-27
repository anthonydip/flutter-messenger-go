package users

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockserver "github.com/anthonydip/flutter-messenger-go/app/storefront-api/webserver/mock"
	mockstore "github.com/anthonydip/flutter-messenger-go/internal/storefront/mock"

	"github.com/gorilla/mux"
)

func TestPost(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		expectedCode     int
		storefrontResult mockstore.Result
	}{
		"created": {
			expectedCode: 201,
		},
		"conflict": {
			expectedCode:     409,
			storefrontResult: mockstore.PostUserResult(errors.New("409 Conflict")),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			srv := mockserver.New().WithStorefront(test.storefrontResult)

			r := mux.NewRouter()
			r.HandleFunc("/users", Post(srv)).Methods(http.MethodPost)

			requestBody := []byte(`{"email": "mock@storefront-mock.com", "provider": "Flutter", "password": "test123123"}`)

			req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatalf("couldn't create test HTTP request: %s", err.Error())
			}

			req.Header.Add("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != test.expectedCode {
				t.Fatalf("expected status code %03d but got %03d (body: %s)", test.expectedCode, rr.Code, rr.Body)
			}
		})
	}
}

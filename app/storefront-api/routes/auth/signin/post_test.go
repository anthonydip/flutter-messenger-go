package signin

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
		"success": {
			expectedCode: 200,
		},
		"conflict": {
			expectedCode:     404,
			storefrontResult: mockstore.SignInResult(errors.New("user does not exist")),
		},
		"unauthorized": {
			expectedCode:     401,
			storefrontResult: mockstore.SignInResult(errors.New("invalid password")),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			srv := mockserver.New().WithStorefront(test.storefrontResult)

			r := mux.NewRouter()
			r.HandleFunc("/auth/signin", Post(srv)).Methods(http.MethodPost)

			requestBody := []byte(`{"email": "mock@storefront-mock.com", "provider": "Flutter", "password": "test123123"}`)

			req, err := http.NewRequest(http.MethodPost, "/auth/signin", bytes.NewBuffer(requestBody))
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

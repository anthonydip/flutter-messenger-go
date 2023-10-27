package users

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockserver "github.com/anthonydip/flutter-messenger-go/app/storefront-api/webserver/mock"
	mockstore "github.com/anthonydip/flutter-messenger-go/internal/storefront/mock"

	"github.com/gorilla/mux"
)

func TestGet(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		expectedCode     int
		storefrontResult mockstore.Result
	}{
		"found": {
			expectedCode: 200,
		},
		"not found": {
			expectedCode:     404,
			storefrontResult: mockstore.GetUserResult(errors.New("user not found")),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			srv := mockserver.New().WithStorefront(test.storefrontResult)

			r := mux.NewRouter()
			r.HandleFunc("/users/{userID}", Get(srv)).Methods(http.MethodGet)

			req, err := http.NewRequest(http.MethodGet, "/users/8ae84a23-fa49-45eb-8000-bdc9b9fe074a", nil)
			if err != nil {
				t.Fatalf("couldn't create test HTTP request: %s", err.Error())
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != test.expectedCode {
				t.Fatalf("expected status code %03d but got %03d (body: %s)", test.expectedCode, rr.Code, rr.Body)
			}
		})
	}
}

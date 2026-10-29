package checker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nxvrlxv/hosts_monitor/internal/config"
)

func TestCheckStatus(t *testing.T) {
	tests := []struct {
		name              string
		handler           http.HandlerFunc
		expectedStatus    int
		expectedCondition bool
	}{
		{
			name: "Status OK",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAccepted)
			},
			expectedStatus:    202,
			expectedCondition: true,
		},
		{
			name: "Status Redirect",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusPermanentRedirect)
			},
			expectedStatus:    308,
			expectedCondition: true,
		},
		{
			name: "Status 500",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectedStatus:    500,
			expectedCondition: false,
		},
	}

	testChecker := Checker{
		Client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		ts := httptest.NewServer(test.handler)
		defer ts.Close()
		testTarget := config.Target{
			ID:  1,
			URL: ts.URL,
		}
		res := testChecker.Check(ctx, testTarget)

		if res.Status != test.expectedStatus {
			t.Errorf("error: expected %v, got %v", test.expectedStatus, res.Status)
		}

		if res.IsUp != test.expectedCondition {
			t.Errorf("error: expected %v, got %v", test.expectedCondition, res.IsUp)
		}

	}
}

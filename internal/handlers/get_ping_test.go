package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPinger struct {
	mock.Mock
}

func (m *MockPinger) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func TestMetricsHandler_GetPing(t *testing.T) {
	tests := []struct {
		name         string
		mockPingErr  error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Ping success",
			mockPingErr:  nil,
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			name:         "Ping failure",
			mockPingErr:  errors.New("ping error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Internal Error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockPinger := new(MockPinger)
			mockPinger.On("Ping").Return(tc.mockPingErr)

			mh := &MetricsHandler{pinger: mockPinger}

			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/ping", nil)

			mh.GetPing(res, req)

			assert.Equal(t, tc.expectedCode, res.Code)

			if tc.expectedBody != "" {
				assert.Contains(t, res.Body.String(), tc.expectedBody)
			}

			mockPinger.AssertExpectations(t)
		})
	}
}

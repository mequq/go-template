package handler

import (
	healthzusecase "application/internal/biz/healthz"
	mocks "application/internal/biz/healthz/mocks"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestHealthzHandler_HealthzLiveness(t *testing.T) {

	type fields struct {
		logger *slog.Logger
		uc     func() *mocks.MockHealthzUseCaseInterface
	}
	type request struct {
		method string
	}

	type response struct {
		status int
	}

	tests := []struct {
		name     string
		fields   fields
		request  request
		response response
	}{
		{
			name: "success",
			fields: fields{
				logger: slog.Default(),

				uc: func() *mocks.MockHealthzUseCaseInterface {
					uc := mocks.NewMockHealthzUseCaseInterface(t)
					// uc.EXPECT().Liveness(nil).Return(nil)
					uc.EXPECT().Liveness(mock.Anything).Return(nil)
					return uc
				},
			},
			request: request{
				method: "GET",
			},
			response: response{
				status: 200,
			},
		},
		{
			name: "error",
			fields: fields{
				logger: slog.Default(),
				uc: func() *mocks.MockHealthzUseCaseInterface {
					uc := mocks.NewMockHealthzUseCaseInterface(t)
					// uc.EXPECT().Liveness(nil).Return(nil)
					uc.EXPECT().Liveness(mock.Anything).Return(errors.New("error"))
					return uc
				},
			},
			request: request{
				method: http.MethodGet,
			},

			response: response{
				status: 500,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := tt.fields.uc()
			s := &HealthzHandler{
				logger: tt.fields.logger,
				uc:     uc,
			}
			recoder := httptest.NewRecorder()
			request := httptest.NewRequest(tt.request.method, "/healthz/liveness", nil) //nolint

			s.HealthzLiveness(recoder, request)

			uc.AssertCalled(t, "Liveness", mock.Anything)
			uc.AssertNotCalled(t, "Readiness", mock.Anything)
			uc.AssertNumberOfCalls(t, "Liveness", 1)
			if recoder.Code != tt.response.status {
				t.Errorf("HealthzHandler.HealthzLiveness() got = %v, want %v", recoder.Code, tt.response.status)
			}

		})
	}
}

func TestNewMuxHealthzHandler(t *testing.T) {
	type args struct {
		uc     func() healthzusecase.HealthzUseCaseInterface
		logger *slog.Logger
	}
	tests := []struct {
		name string
		args args
		want *HealthzHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMuxHealthzHandler(tt.args.uc(), tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMuxHealthzHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHealthzHandler_HealthzReadiness(t *testing.T) {

	type fields struct {
		logger *slog.Logger
		mocker func() *mocks.MockHealthzUseCaseInterface
	}

	type request struct {
		method string
	}
	type response struct {
		status int
	}

	tests := []struct {
		name     string
		fields   fields
		request  request
		response response
	}{
		{
			name: "success",
			request: request{
				method: "GET",
			},

			response: response{
				status: 200,
			},
			fields: fields{
				logger: slog.Default(),
				mocker: func() *mocks.MockHealthzUseCaseInterface {
					uc := mocks.NewMockHealthzUseCaseInterface(t)
					uc.EXPECT().Readiness(mock.Anything).Return(nil)
					return uc
				},
			},
		},
		{
			name: "error",
			request: request{
				method: http.MethodGet,
			},

			response: response{
				status: 500,
			},

			fields: fields{
				logger: slog.Default(),
				mocker: func() *mocks.MockHealthzUseCaseInterface {
					uc := mocks.NewMockHealthzUseCaseInterface(t)
					uc.EXPECT().Readiness(mock.Anything).Return(errors.New("error"))
					return uc
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.fields.mocker()
			s := &HealthzHandler{
				logger: tt.fields.logger,
				uc:     m,
			}

			recoder := httptest.NewRecorder()
			request := httptest.NewRequest(tt.request.method, "/healthz/readiness", nil) //nolint

			s.HealthzReadiness(recoder, request)

			m.AssertCalled(t, "Readiness", mock.Anything)
			m.AssertNotCalled(t, "Liveness", mock.Anything)
			m.AssertNumberOfCalls(t, "Readiness", 1)

			if recoder.Code != tt.response.status {
				t.Errorf("HealthzHandler.HealthzReadiness() got = %v, want %v", recoder.Code, tt.response.status)
			}
		})
	}
}

func TestHealthzHandler_Panic(t *testing.T) {
	type fields struct {
		logger *slog.Logger
		uc     healthzusecase.HealthzUseCaseInterface
	}
	type args struct {
		in0 http.ResponseWriter
		in1 *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HealthzHandler{
				logger: tt.fields.logger,
				uc:     tt.fields.uc,
			}
			s.Panic(tt.args.in0, tt.args.in1)
		})
	}
}

func TestHealthzHandler_LongRun(t *testing.T) {
	type fields struct {
		logger *slog.Logger
		uc     healthzusecase.HealthzUseCaseInterface
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HealthzHandler{
				logger: tt.fields.logger,
				uc:     tt.fields.uc,
			}
			s.LongRun(tt.args.w, tt.args.r)
		})
	}
}

func TestHealthzHandler_RegisterMuxRouter(t *testing.T) {
	type fields struct {
		logger *slog.Logger
		uc     healthzusecase.HealthzUseCaseInterface
	}
	type args struct {
		mux *http.ServeMux
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HealthzHandler{
				logger: tt.fields.logger,
				uc:     tt.fields.uc,
			}
			s.RegisterMuxRouter(tt.args.mux)
		})
	}
}

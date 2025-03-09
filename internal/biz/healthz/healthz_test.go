package healthzusecase

import (
	mockhealthzrepo "application/internal/biz/healthz/mocks"
	"context"
	"log/slog"
	"testing"

	"errors"

	"github.com/stretchr/testify/mock"
	tracenoopotel "go.opentelemetry.io/otel/trace/noop"
)

func TestHealthzUseCase_Liveness(t *testing.T) {

	type fields struct {
		repo   func() HealthzRepoInterface
		logger *slog.Logger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test01",
			fields: fields{
				repo: func() HealthzRepoInterface {
					m := mockhealthzrepo.NewMockHealthzRepoInterface(t)
					m.EXPECT().Liveness(mock.Anything).Return(nil)
					return m

				},
				logger: slog.Default(),
			},
			args:    args{ctx: context.Background()},
			wantErr: false,
		},
		{
			name: "test02",
			fields: fields{
				repo: func() HealthzRepoInterface {
					m := mockhealthzrepo.NewMockHealthzRepoInterface(t)
					m.EXPECT().Liveness(mock.Anything).Return(errors.New("error"))
					return m

				},
				logger: slog.Default(),
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &HealthzBiz{
				repo:   tt.fields.repo(),
				logger: tt.fields.logger,
				tracer: tracenoopotel.NewTracerProvider().Tracer(""),
			}
			if err := uc.Liveness(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("HealthzBiz.Liveness() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package healthzusecase

import (
	"context"
	"log/slog"
	"reflect"
	"testing"
)

type mockRepo struct{}

func NewRepoMock() HealthzRepoInterface {
	return &mockRepo{}
}

func (m *mockRepo) Readiness(_ context.Context) error {
	return nil
}

func (m *mockRepo) Liveness(_ context.Context) error {
	return nil
}

func TestHealthzUseCase_Liveness(t *testing.T) {
	type fields struct {
		repo   HealthzRepoInterface
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
				repo:   NewRepoMock(),
				logger: slog.Default(),
			},
			args:    args{ctx: context.Background()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &HealthzBiz{
				repo:   tt.fields.repo,
				logger: tt.fields.logger,
			}
			if err := uc.Liveness(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("HealthzBiz.Liveness() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHealthzUseCase_Readiness(t *testing.T) {
	type fields struct {
		repo   HealthzRepoInterface
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
			name: "readiness-test",
			fields: fields{
				repo:   NewRepoMock(),
				logger: slog.Default(),
			},
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &HealthzBiz{
				repo:   tt.fields.repo,
				logger: tt.fields.logger,
			}
			if err := uc.Readiness(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("HealthzBiz.Readiness() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewHealthzUseCase(t *testing.T) {
	logger := slog.Default()
	huc := &HealthzBiz{
		repo:   nil,
		logger: logger,
	}

	type args struct {
		repo   HealthzRepoInterface
		logger *slog.Logger
	}

	tests := []struct {
		name string
		args args
		want HealthzUseCaseInterface
	}{
		{
			name: "newHealthzUsecase",
			args: args{
				repo:   nil,
				logger: slog.Default(),
			},
			want: huc,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHealthzBiz(tt.args.repo, tt.args.logger); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("NewHealthzBiz() = %v, want %v", got, tt.want)
			}
		})
	}
}

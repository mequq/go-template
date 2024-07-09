package sample_entity

import (
	"application/internal/datasource/sample_entitiy"
	"application/internal/entity"
	mse "application/mocks/datasource"
	"context"
	"errors"
	"go.uber.org/mock/gomock"
	"log/slog"
	"os"
	"testing"
)

func TestUserService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	var tests = []struct {
		name               string
		sampleEntityDSMock func() *mse.MockDataSource
		seInput            *entity.SampleEntity
		seOutput           *entity.SampleEntity
		ctx                context.Context
		error              error
	}{
		{
			name: "success",
			sampleEntityDSMock: func() *mse.MockDataSource {
				dsMock := mse.NewMockDataSource(ctrl)
				dsMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uint64(1), nil)
				return dsMock
			},
			seInput: &entity.SampleEntity{
				ID:   0,
				Name: "name",
				Text: "text",
			},
			seOutput: &entity.SampleEntity{
				ID:   1,
				Name: "name",
				Text: "text",
			},
			ctx:   context.Background(),
			error: nil,
		},
		{
			name: "already-exist",
			sampleEntityDSMock: func() *mse.MockDataSource {
				dsMock := mse.NewMockDataSource(ctrl)
				dsMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uint64(1), sample_entitiy.ErrAlreadyExist)
				return dsMock
			},
			seInput: &entity.SampleEntity{
				ID:   0,
				Name: "name",
				Text: "text",
			},
			seOutput: nil,
			ctx:      context.Background(),
			error:    sample_entitiy.ErrAlreadyExist,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			seDSMock := test.sampleEntityDSMock()
			biz := NewSampleEntity(seDSMock, slog.New(slog.NewTextHandler(os.Stdout,
				nil)))
			se, err := biz.Create(test.ctx, test.seInput)
			if !errors.Is(err, test.error) {
				t.Errorf("error:%s is not equal to %s", err, test.error)
			}
			if !gomock.Eq(se).Matches(test.seOutput) {
				t.Errorf("output:%v is not equal to %v", se, test.seOutput)
			}

			seDSMock.EXPECT()
		})
	}
}

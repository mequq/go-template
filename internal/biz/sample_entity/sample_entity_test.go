package sample_entity

import (
	"application/internal/datasource/sample_entitiy"
	"application/internal/entity"
	mock_sample_entitiy "application/mocks/datasource"
	"context"
	"errors"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestUserService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})
	err := errors.New("error")

	var tests = []struct {
		name               string
		sampleEntityDSMock func() *mock_sample_entitiy.MockDataSource
		seInput            *entity.SampleEntity
		seOutput           *entity.SampleEntity
		ctx                context.Context
		error              error
	}{
		{
			name: "success",
			sampleEntityDSMock: func() *mock_sample_entitiy.MockDataSource {
				userRepoMock := mock_sample_entitiy.NewMockDataSource(ctrl)
				userRepoMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				return userRepoMock
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
			name: "not-found",
			sampleEntityDSMock: func() *mock_sample_entitiy.MockDataSource {
				userRepoMock := mock_sample_entitiy.NewMockDataSource(ctrl)
				userRepoMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				return userRepoMock
			},
			seInput: &entity.SampleEntity{
				ID:   0,
				Name: "name",
				Text: "text",
			},
			seOutput: nil,
			ctx:      context.Background(),
			error:    sample_entitiy.ErrNotFound,
		},
		{
			name: "already-exist",
			sampleEntityDSMock: func() *mock_sample_entitiy.MockDataSource {
				userRepoMock := mock_sample_entitiy.NewMockDataSource(ctrl)
				userRepoMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				return userRepoMock
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
			userRepoMock := test.userRepoMock()
			loggerMock := test.loggerMock()
			service := NewUserService(userRepoMock, loggerMock)
			err := service.CreateUser(test.user)

			if !errors.Is(err, test.error) {
				t.Errorf("error:%s is not equal to %s", err, test.error)
			}
			loggerMock.EXPECT()
			userRepoMock.EXPECT()
		})
	}
}

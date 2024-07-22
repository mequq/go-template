package sample_entity

import (
	"application/internal/v1/datasource/sample_entitiy"
	"application/internal/v1/entity"
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	mse "application/mock/datasource"

	"go.uber.org/mock/gomock"
)

func TestUserService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	tests := []struct {
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

func BenchmarkUserService_CreateUser(b *testing.B) {
	ctrl := gomock.NewController(b)
	dsMock := mse.NewMockDataSource(ctrl)
	dsMock.EXPECT().Create(gomock.Any(), gomock.Any()).Times(b.N).Return(uint64(1), nil).AnyTimes()
	biz := NewSampleEntity(dsMock, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	se := &entity.SampleEntity{
		ID:   0,
		Name: "name",
		Text: "text",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		biz.Create(context.Background(), se)
	}
}

func TestSampleEntity_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	dbErr := errors.New("database error")
	tests := []struct {
		name               string
		sampleEntityDSMock func() *mse.MockDataSource
		expectedOutput     []*entity.SampleEntity
		ctx                context.Context
		error              error
	}{
		{
			name: "success",
			sampleEntityDSMock: func() *mse.MockDataSource {
				dsMock := mse.NewMockDataSource(ctrl)
				dsMock.EXPECT().List(gomock.Any()).Return([]*entity.SampleEntity{
					{
						ID:   1,
						Name: "name1",
						Text: "text1",
					},
					{
						ID:   2,
						Name: "name2",
						Text: "text2",
					},
				}, nil)
				return dsMock
			},
			expectedOutput: []*entity.SampleEntity{
				{
					ID:   1,
					Name: "name1",
					Text: "text1",
				},
				{
					ID:   2,
					Name: "name2",
					Text: "text2",
				},
			},
			ctx:   context.Background(),
			error: nil,
		},
		{
			name: "error",
			sampleEntityDSMock: func() *mse.MockDataSource {
				dsMock := mse.NewMockDataSource(ctrl)
				dsMock.EXPECT().List(gomock.Any()).Return(nil, dbErr)
				return dsMock
			},
			expectedOutput: nil,
			ctx:            context.Background(),
			error:          dbErr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			seDSMock := test.sampleEntityDSMock()
			biz := NewSampleEntity(seDSMock, slog.New(slog.NewTextHandler(os.Stdout, nil)))
			se, err := biz.List(test.ctx)
			if !errors.Is(err, test.error) {
				t.Errorf("error:%s is not equal to %s", err, test.error)
			}
			if !gomock.Eq(se).Matches(test.expectedOutput) {
				t.Errorf("output:%v is not equal to %v", se, test.expectedOutput)
			}
		})
	}
}

func BenchmarkSampleEntity_List(b *testing.B) {
	ctrl := gomock.NewController(b)
	dsMock := mse.NewMockDataSource(ctrl)
	dsMock.EXPECT().List(gomock.Any()).Times(b.N).Return([]*entity.SampleEntity{
		{
			ID:   1,
			Name: "name1",
			Text: "text1",
		},
		{
			ID:   2,
			Name: "name2",
			Text: "text2",
		},
	}, nil).AnyTimes()
	biz := NewSampleEntity(dsMock, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		biz.List(context.Background())
	}
}

func TestSampleEntity_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	dbErr := errors.New("database error")
	tests := []struct {
		name               string
		sampleEntityDSMock func() *mse.MockDataSource
		seInput            *entity.SampleEntity
		ctx                context.Context
		error              error
	}{
		{
			name: "success",
			sampleEntityDSMock: func() *mse.MockDataSource {
				dsMock := mse.NewMockDataSource(ctrl)
				dsMock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				return dsMock
			},
			seInput: &entity.SampleEntity{
				ID:   1,
				Name: "name",
				Text: "text",
			},
			ctx:   context.Background(),
			error: nil,
		},
		{
			name: "not-found",
			sampleEntityDSMock: func() *mse.MockDataSource {
				dsMock := mse.NewMockDataSource(ctrl)
				dsMock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(sample_entitiy.ErrNotFound)
				return dsMock
			},
			seInput: &entity.SampleEntity{
				ID:   1,
				Name: "name",
				Text: "text",
			},
			ctx:   context.Background(),
			error: sample_entitiy.ErrNotFound,
		},
		{
			name: "already-exist",
			sampleEntityDSMock: func() *mse.MockDataSource {
				dsMock := mse.NewMockDataSource(ctrl)
				dsMock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(sample_entitiy.ErrAlreadyExist)
				return dsMock
			},
			seInput: &entity.SampleEntity{
				ID:   1,
				Name: "name",
				Text: "text",
			},
			ctx:   context.Background(),
			error: sample_entitiy.ErrAlreadyExist,
		},
		{
			name: "error",
			sampleEntityDSMock: func() *mse.MockDataSource {
				dsMock := mse.NewMockDataSource(ctrl)
				dsMock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(dbErr)
				return dsMock
			},
			seInput: &entity.SampleEntity{
				ID:   1,
				Name: "name",
				Text: "text",
			},
			ctx:   context.Background(),
			error: dbErr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			seDSMock := test.sampleEntityDSMock()
			biz := NewSampleEntity(seDSMock, slog.New(slog.NewTextHandler(os.Stdout, nil)))
			err := biz.Update(test.ctx, test.seInput)
			if !errors.Is(err, test.error) {
				t.Errorf("error:%s is not equal to %s", err, test.error)
			}
		})
	}
}

func BenchmarkSampleEntity_Update(b *testing.B) {
	ctrl := gomock.NewController(b)
	dsMock := mse.NewMockDataSource(ctrl)
	dsMock.EXPECT().Update(gomock.Any(), gomock.Any()).Times(b.N).Return(nil).AnyTimes()
	biz := NewSampleEntity(dsMock, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	se := &entity.SampleEntity{
		ID:   1,
		Name: "name",
		Text: "text",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		biz.Update(context.Background(), se)
	}
}

func TestSampleEntity_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	dbErr := errors.New("database error")
	tests := []struct {
		name               string
		sampleEntityDSMock func() *mse.MockDataSource
		id                 uint64
		ctx                context.Context
		error              error
	}{
		{
			name: "success",
			sampleEntityDSMock: func() *mse.MockDataSource {
				dsMock := mse.NewMockDataSource(ctrl)
				dsMock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				return dsMock
			},
			id:    1,
			ctx:   context.Background(),
			error: nil,
		},
		{
			name: "not-found",
			sampleEntityDSMock: func() *mse.MockDataSource {
				dsMock := mse.NewMockDataSource(ctrl)
				dsMock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(sample_entitiy.ErrNotFound)
				return dsMock
			},
			id:    1,
			ctx:   context.Background(),
			error: sample_entitiy.ErrNotFound,
		},
		{
			name: "error",
			sampleEntityDSMock: func() *mse.MockDataSource {
				dsMock := mse.NewMockDataSource(ctrl)
				dsMock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(dbErr)
				return dsMock
			},
			id:    1,
			ctx:   context.Background(),
			error: dbErr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			seDSMock := test.sampleEntityDSMock()
			biz := NewSampleEntity(seDSMock, slog.New(slog.NewTextHandler(os.Stdout, nil)))
			err := biz.Delete(test.ctx, test.id)
			if !errors.Is(err, test.error) {
				t.Errorf("error:%s is not equal to %s", err, test.error)
			}
		})
	}
}

func BenchmarkSampleEntity_Delete(b *testing.B) {
	ctrl := gomock.NewController(b)
	dsMock := mse.NewMockDataSource(ctrl)
	dsMock.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(b.N).Return(nil).AnyTimes()
	biz := NewSampleEntity(dsMock, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		biz.Delete(context.Background(), 1)
	}
}

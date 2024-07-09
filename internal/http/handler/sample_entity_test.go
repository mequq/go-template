package handler

import (
	"application/internal/datasource/sample_entitiy"
	"application/internal/entity"
	"application/internal/http/dto"
	apiResponse "application/internal/http/response"
	mockBiz "application/mock/biz"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/mock/gomock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSampleEntitieHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	var tests = []struct {
		name                string
		sampleEntityBizMock func() *mockBiz.MockSampleEntity
		request             func() *http.Request
		expectedResponse    apiResponse.Response[[]dto.SampleEntityResponse]
		expectedStatusCode  int
		ctx                 context.Context
	}{
		{
			name: "success",
			sampleEntityBizMock: func() *mockBiz.MockSampleEntity {
				dsMock := mockBiz.NewMockSampleEntity(ctrl)
				dsMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&entity.SampleEntity{
					ID:   1,
					Name: "name",
					Text: "text",
				}, nil)
				return dsMock
			},
			request: func() *http.Request {
				sampleReq := dto.SampleEntityRequest{
					Name: "name",
					Text: "text",
				}
				b, _ := json.Marshal(sampleReq)
				r := bytes.NewReader(b)
				return httptest.NewRequest(http.MethodPost, "/sample-entities", r)
			},
			expectedResponse: apiResponse.Response[[]dto.SampleEntityResponse]{
				Message: "Created Successfully",
				Status:  http.StatusCreated,
				Data:    nil,
			},
			ctx:                context.Background(),
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "already-exist",
			sampleEntityBizMock: func() *mockBiz.MockSampleEntity {
				dsMock := mockBiz.NewMockSampleEntity(ctrl)
				dsMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, sample_entitiy.ErrAlreadyExist)
				return dsMock
			},
			request: func() *http.Request {
				sampleReq := dto.SampleEntityRequest{
					Name: "name",
					Text: "text",
				}
				b, _ := json.Marshal(sampleReq)
				r := bytes.NewReader(b)
				return httptest.NewRequest(http.MethodPost, "/sample-entities", r)
			},
			expectedResponse: apiResponse.Response[[]dto.SampleEntityResponse]{
				Message: "already-exist",
				Status:  http.StatusBadRequest,
				Data:    nil,
			},
			ctx:                context.Background(),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "internal server error",
			sampleEntityBizMock: func() *mockBiz.MockSampleEntity {
				dsMock := mockBiz.NewMockSampleEntity(ctrl)
				dsMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				return dsMock
			},
			request: func() *http.Request {
				sampleReq := dto.SampleEntityRequest{
					Name: "name",
					Text: "text",
				}
				b, _ := json.Marshal(sampleReq)
				r := bytes.NewReader(b)
				return httptest.NewRequest(http.MethodPost, "/sample-entities", r)
			},
			expectedResponse: apiResponse.Response[[]dto.SampleEntityResponse]{
				Message: "internal-error",
				Status:  http.StatusInternalServerError,
				Data:    nil,
			},
			ctx:                context.Background(),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bizMock := test.sampleEntityBizMock()
			handler := NewSampleEntityHandler(slog.New(slog.NewTextHandler(os.Stdout, nil)), bizMock)
			recorder := httptest.NewRecorder()
			handler.Create(recorder, test.request())

			if recorder.Code != test.expectedStatusCode {
				t.Errorf("status code:%d did not match expected value:%d", recorder.Code, test.expectedStatusCode)
			}
			decoder := json.NewDecoder(recorder.Body)
			r := apiResponse.Response[[]dto.SampleEntityResponse]{}
			if err := decoder.Decode(&r); err != nil {
				t.Error(err)
			}
			if !gomock.Eq(r).Matches(test.expectedResponse) {
				t.Errorf("response body not match have:%v want:%v", r, test.expectedResponse)
			}
			bizMock.EXPECT()
		})
	}
}
func BenchmarkSampleEntity_Create(b *testing.B) {
	ctrl := gomock.NewController(b)
	seBiz := mockBiz.NewMockSampleEntity(ctrl)
	seBiz.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&entity.SampleEntity{
		ID:   1,
		Name: "name",
		Text: "text",
	}, nil).AnyTimes()

	sampleReq := dto.SampleEntityRequest{
		Name: "name",
		Text: "text",
	}
	bs, _ := json.Marshal(sampleReq)
	r := bytes.NewReader(bs)

	handler := NewSampleEntityHandler(slog.New(slog.NewTextHandler(os.Stdout, nil)), seBiz)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.Create(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/sample-entities", r))
	}

}

func TestSampleEntitieHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	var tests = []struct {
		name                string
		sampleEntityBizMock func() *mockBiz.MockSampleEntity
		request             func() *http.Request
		expectedResponse    apiResponse.Response[[]dto.SampleEntityResponse]
		expectedStatusCode  int
		ctx                 context.Context
	}{
		{
			name: "success",
			sampleEntityBizMock: func() *mockBiz.MockSampleEntity {
				dsMock := mockBiz.NewMockSampleEntity(ctrl)
				dsMock.EXPECT().List(gomock.Any()).Return([]*entity.SampleEntity{
					{ID: 1, Name: "name1", Text: "text1"}, {ID: 2, Name: "name2", Text: "text2"},
				}, nil)
				return dsMock
			},
			request: func() *http.Request {
				sampleReq := dto.SampleEntityRequest{
					Name: "name",
					Text: "text",
				}
				b, _ := json.Marshal(sampleReq)
				r := bytes.NewReader(b)
				return httptest.NewRequest(http.MethodPost, "/sample-entities", r)
			},
			expectedResponse: apiResponse.Response[[]dto.SampleEntityResponse]{
				Message: "",
				Status:  http.StatusOK,
				Data: dto.SampleEntityListResponses([]*entity.SampleEntity{
					{ID: 1, Name: "name1", Text: "text1"}, {ID: 2, Name: "name2", Text: "text2"},
				}),
			},
			ctx:                context.Background(),
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "internal server error",
			sampleEntityBizMock: func() *mockBiz.MockSampleEntity {
				dsMock := mockBiz.NewMockSampleEntity(ctrl)
				dsMock.EXPECT().List(gomock.Any()).Return(nil, errors.New("error"))
				return dsMock
			},
			request: func() *http.Request {
				r := bytes.NewReader([]byte{})
				return httptest.NewRequest(http.MethodGet, "/sample-entities", r)
			},
			expectedResponse: apiResponse.Response[[]dto.SampleEntityResponse]{
				Message: "internal-error",
				Status:  http.StatusInternalServerError,
				Data:    nil,
			},
			ctx:                context.Background(),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bizMock := test.sampleEntityBizMock()
			handler := NewSampleEntityHandler(slog.New(slog.NewTextHandler(os.Stdout, nil)), bizMock)
			recorder := httptest.NewRecorder()
			handler.List(recorder, test.request())

			if recorder.Code != test.expectedStatusCode {
				t.Errorf("status code:%d did not match expected value:%d", recorder.Code, test.expectedStatusCode)
			}
			decoder := json.NewDecoder(recorder.Body)
			r := apiResponse.Response[[]dto.SampleEntityResponse]{}
			if err := decoder.Decode(&r); err != nil {
				t.Error(err)
			}
			if !gomock.Eq(r).Matches(test.expectedResponse) {
				fmt.Println(r, test.expectedResponse)
				t.Errorf("response body not match have:%v want:%v", r, test.expectedResponse)
			}
			bizMock.EXPECT()
		})
	}
}

func BenchmarkSampleEntity_List(b *testing.B) {
	ctrl := gomock.NewController(b)
	seBiz := mockBiz.NewMockSampleEntity(ctrl)
	seBiz.EXPECT().List(gomock.Any()).Return([]*entity.SampleEntity{
		{ID: 1, Name: "name1", Text: "text1"}, {ID: 2, Name: "name2", Text: "text2"},
	}, nil).AnyTimes()

	sampleReq := dto.SampleEntityRequest{
		Name: "name",
		Text: "text",
	}
	bs, _ := json.Marshal(sampleReq)
	r := bytes.NewReader(bs)

	handler := NewSampleEntityHandler(slog.New(slog.NewTextHandler(os.Stdout, nil)), seBiz)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.List(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/sample-entities", r))
	}
}

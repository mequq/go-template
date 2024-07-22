package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"application/internal/v1/datasource/sampleentity"
	"application/internal/v1/entity"
	"application/internal/v1/http/dto"

	apiResponse "application/internal/v1/http/response"
	mockBiz "application/mock/biz"
	"go.uber.org/mock/gomock"
)

func TestSampleEntitieHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	tests := []struct {
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
				dsMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, sampleentity.ErrAlreadyExist)
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

	tests := []struct {
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

func TestSampleEntityHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	tests := []struct {
		name                string
		sampleEntityBizMock func() *mockBiz.MockSampleEntity
		request             func() *http.Request
		expectedStatusCode  int
		expectedResponse    apiResponse.Response[any]
	}{
		{
			name: "success",
			sampleEntityBizMock: func() *mockBiz.MockSampleEntity {
				dsMock := mockBiz.NewMockSampleEntity(ctrl)
				dsMock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				return dsMock
			},
			request: func() *http.Request {
				sampleReq := dto.SampleEntityRequest{
					Name: "updatedName",
					Text: "updatedText",
				}
				b, _ := json.Marshal(sampleReq)
				r := httptest.NewRequest(http.MethodPut, "/entities/1/", bytes.NewReader(b))
				r.SetPathValue("id", "1")
				return r
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: apiResponse.Response[any]{
				Message: "Updated successfully",
				Status:  http.StatusOK,
				Data:    nil,
			},
		},
		{
			name: "invalid request",
			sampleEntityBizMock: func() *mockBiz.MockSampleEntity {
				return mockBiz.NewMockSampleEntity(ctrl)
			},
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodPut, "/entities/abc", nil)
				return r
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: apiResponse.Response[any]{
				Message: "invalid-request",
				Status:  http.StatusBadRequest,
				Data:    nil,
			},
		},
		{
			name: "not found",
			sampleEntityBizMock: func() *mockBiz.MockSampleEntity {
				dsMock := mockBiz.NewMockSampleEntity(ctrl)
				dsMock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(sampleentity.ErrNotFound)
				return dsMock
			},
			request: func() *http.Request {
				sampleReq := dto.SampleEntityRequest{
					Name: "updatedName",
					Text: "updatedText",
				}
				b, _ := json.Marshal(sampleReq)
				r := httptest.NewRequest(http.MethodPut, "/entities/1/", bytes.NewReader(b))
				r.SetPathValue("id", "1")
				return r
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: apiResponse.Response[any]{
				Message: "not-found",
				Status:  http.StatusNotFound,
				Data:    nil,
			},
		},
		{
			name: "internal server error",
			sampleEntityBizMock: func() *mockBiz.MockSampleEntity {
				dsMock := mockBiz.NewMockSampleEntity(ctrl)
				dsMock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("database error"))
				return dsMock
			},
			request: func() *http.Request {
				sampleReq := dto.SampleEntityRequest{
					Name: "updatedName",
					Text: "updatedText",
				}
				b, _ := json.Marshal(sampleReq)
				r := httptest.NewRequest(http.MethodPut, "/entities/1", bytes.NewReader(b))
				r.SetPathValue("id", "1")
				return r
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: apiResponse.Response[any]{
				Message: "internal-error",
				Status:  http.StatusInternalServerError,
				Data:    nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bizMock := test.sampleEntityBizMock()
			handler := NewSampleEntityHandler(slog.New(slog.NewTextHandler(os.Stdout, nil)), bizMock)
			recorder := httptest.NewRecorder()
			handler.Update(recorder, test.request())

			if recorder.Code != test.expectedStatusCode {
				t.Errorf("status code:%d did not match expected value:%d", recorder.Code, test.expectedStatusCode)
			}

			var r apiResponse.Response[any]
			if err := json.NewDecoder(recorder.Body).Decode(&r); err != nil {
				t.Errorf("error decoding response body: %v", err)
			}

			if !gomock.Eq(r).Matches(test.expectedResponse) {
				t.Errorf("response body not match have:%v want:%v", r, test.expectedResponse)
			}
		})
	}
}

func BenchmarkSampleEntityHandler_Update(b *testing.B) {
	ctrl := gomock.NewController(b)
	seBiz := mockBiz.NewMockSampleEntity(ctrl)
	seBiz.EXPECT().Update(gomock.Any(), gomock.Any()).Times(b.N).Return(nil).AnyTimes()

	handler := NewSampleEntityHandler(slog.New(slog.NewTextHandler(os.Stdout, nil)), seBiz)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		recorder := httptest.NewRecorder()
		sampleReq := dto.SampleEntityRequest{
			Name: "updatedName",
			Text: "updatedText",
		}
		b, _ := json.Marshal(sampleReq)
		req := httptest.NewRequest(http.MethodPut, "/entities/1", bytes.NewReader(b))
		handler.Update(recorder, req)
	}
}

func TestSampleEntityHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	tests := []struct {
		name                string
		sampleEntityBizMock func() *mockBiz.MockSampleEntity
		request             func() *http.Request
		expectedStatusCode  int
		expectedResponse    apiResponse.Response[any]
	}{
		{
			name: "success",
			sampleEntityBizMock: func() *mockBiz.MockSampleEntity {
				dsMock := mockBiz.NewMockSampleEntity(ctrl)
				dsMock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				return dsMock
			},
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodDelete, "/entities/1/", bytes.NewReader([]byte{}))
				r.SetPathValue("id", "1")
				return r
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: apiResponse.Response[any]{
				Message: "sample entity deleted",
				Status:  http.StatusOK,
				Data:    nil,
			},
		},
		{
			name: "invalid id",
			sampleEntityBizMock: func() *mockBiz.MockSampleEntity {
				return mockBiz.NewMockSampleEntity(ctrl)
			},
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodDelete, "/entities/abc", nil)
				return r
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: apiResponse.Response[any]{
				Message: "invalid-request",
				Status:  http.StatusBadRequest,
				Data:    nil,
			},
		},
		{
			name: "not found",
			sampleEntityBizMock: func() *mockBiz.MockSampleEntity {
				dsMock := mockBiz.NewMockSampleEntity(ctrl)
				dsMock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(sampleentity.ErrNotFound)
				return dsMock
			},
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodDelete, "/entities/1/", bytes.NewReader([]byte{}))
				r.SetPathValue("id", "1")
				return r
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: apiResponse.Response[any]{
				Message: "not-found",
				Status:  http.StatusNotFound,
				Data:    nil,
			},
		},
		{
			name: "internal server error",
			sampleEntityBizMock: func() *mockBiz.MockSampleEntity {
				dsMock := mockBiz.NewMockSampleEntity(ctrl)
				dsMock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("database error"))
				return dsMock
			},
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodDelete, "/entities/1/", bytes.NewReader([]byte{}))
				r.SetPathValue("id", "1")
				return r
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: apiResponse.Response[any]{
				Message: "internal-error",
				Status:  http.StatusInternalServerError,
				Data:    nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bizMock := test.sampleEntityBizMock()
			handler := NewSampleEntityHandler(slog.New(slog.NewTextHandler(os.Stdout, nil)), bizMock)
			recorder := httptest.NewRecorder()
			handler.Delete(recorder, test.request())

			if recorder.Code != test.expectedStatusCode {
				t.Errorf("status code:%d did not match expected value:%d", recorder.Code, test.expectedStatusCode)
			}

			var r apiResponse.Response[any]
			if err := json.NewDecoder(recorder.Body).Decode(&r); err != nil {
				t.Errorf("error decoding response body: %v", err)
			}

			if !gomock.Eq(r).Matches(test.expectedResponse) {
				t.Errorf("response body not match have:%v want:%v", r, test.expectedResponse)
			}
		})
	}
}

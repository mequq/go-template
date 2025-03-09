package handler

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"log/slog"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"

// 	sampleusecasev1 "application/internal/biz/sample"
// 	sampleentity "application/internal/entity/sample"
// 	"application/internal/http/dto"
// 	"application/internal/http/response"
// 	mock_sampleusecasev1 "application/mock/sample/v1"

// 	"go.uber.org/mock/gomock"
// )

// func BenchmarkSampleEntity_Create(b *testing.B) {
// 	ctrl := gomock.NewController(b)
// 	seBiz := mock_sampleusecasev1.NewMockSampleEntityUsecaseInterface(ctrl)
// 	seBiz.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&sampleentity.Sample{
// 		ID:   1,
// 		Name: "name",
// 		Text: "text",
// 	}, nil).AnyTimes()

// 	sampleReq := dto.SampleEntityRequest{
// 		Name: "name",
// 		Text: "text",
// 	}
// 	bs, _ := json.Marshal(sampleReq)
// 	r := bytes.NewReader(bs)

// 	handler := NewSampleEntityHandler(slog.New(slog.NewTextHandler(os.Stdout, nil)), seBiz)

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		handler.Create(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/sample-entities", r))
// 	}
// }

// func TestSampleEntitieHandler_Create(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	t.Cleanup(func() {
// 		ctrl.Finish()
// 	})

// 	tests := []struct {
// 		name                string
// 		sampleEntityBizMock func() *mock_sampleusecasev1.MockSampleEntityUsecaseInterface
// 		request             func() *http.Request
// 		expectedResponse    response.Response[*dto.SampleEntityResponse]
// 		expectedStatusCode  int
// 		ctx                 context.Context
// 	}{
// 		{
// 			name: "success",
// 			sampleEntityBizMock: func() *mock_sampleusecasev1.MockSampleEntityUsecaseInterface {
// 				dsMock := mock_sampleusecasev1.NewMockSampleEntityUsecaseInterface(ctrl)
// 				dsMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&sampleentity.Sample{
// 					ID:   1,
// 					Name: "name",
// 					Text: "text",
// 				}, nil)
// 				return dsMock
// 			},
// 			request: func() *http.Request {
// 				sampleReq := dto.SampleEntityRequest{
// 					Name: "name",
// 					Text: "text",
// 				}
// 				b, _ := json.Marshal(sampleReq)
// 				r := bytes.NewReader(b)
// 				return httptest.NewRequest(http.MethodPost, "/sample-entities", r)
// 			},
// 			expectedResponse: response.Response[*dto.SampleEntityResponse]{
// 				Message: "Created Successfully",
// 				Status:  http.StatusCreated,
// 				Data: &dto.SampleEntityResponse{
// 					ID:   1,
// 					Name: "name",
// 					Text: "text",
// 				},
// 			},
// 			ctx:                context.Background(),
// 			expectedStatusCode: http.StatusCreated,
// 		},
// 		{
// 			name: "already-exist",
// 			sampleEntityBizMock: func() *mock_sampleusecasev1.MockSampleEntityUsecaseInterface {
// 				dsMock := mock_sampleusecasev1.NewMockSampleEntityUsecaseInterface(ctrl)
// 				dsMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, sampleusecasev1.ErrAlreadyExist)
// 				return dsMock
// 			},
// 			request: func() *http.Request {
// 				sampleReq := dto.SampleEntityRequest{
// 					Name: "name",
// 					Text: "text",
// 				}
// 				b, _ := json.Marshal(sampleReq)
// 				r := bytes.NewReader(b)
// 				return httptest.NewRequest(http.MethodPost, "/sample-entities", r)
// 			},
// 			expectedResponse: response.Response[*dto.SampleEntityResponse]{
// 				Message: "already-exist",
// 				Status:  http.StatusBadRequest,
// 				Data:    nil,
// 			},
// 			ctx:                context.Background(),
// 			expectedStatusCode: http.StatusBadRequest,
// 		},
// 		{
// 			name: "internal server error",
// 			sampleEntityBizMock: func() *mock_sampleusecasev1.MockSampleEntityUsecaseInterface {
// 				dsMock := mock_sampleusecasev1.NewMockSampleEntityUsecaseInterface(ctrl)
// 				dsMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
// 				return dsMock
// 			},
// 			request: func() *http.Request {
// 				sampleReq := dto.SampleEntityRequest{
// 					Name: "name",
// 					Text: "text",
// 				}
// 				b, _ := json.Marshal(sampleReq)
// 				r := bytes.NewReader(b)
// 				return httptest.NewRequest(http.MethodPost, "/sample-entities", r)
// 			},
// 			expectedResponse: response.Response[*dto.SampleEntityResponse]{
// 				Message: "internal-error",
// 				Status:  http.StatusInternalServerError,
// 				Data:    nil,
// 			},
// 			ctx:                context.Background(),
// 			expectedStatusCode: http.StatusInternalServerError,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			bizMock := test.sampleEntityBizMock()
// 			handler := NewSampleEntityHandler(slog.New(slog.NewTextHandler(os.Stdout, nil)), bizMock)
// 			recorder := httptest.NewRecorder()
// 			handler.Create(recorder, test.request())

// 			if recorder.Code != test.expectedStatusCode {
// 				t.Errorf("status code:%d did not match expected value:%d", recorder.Code, test.expectedStatusCode)
// 			}
// 			decoder := json.NewDecoder(recorder.Body)
// 			var r response.Response[*dto.SampleEntityResponse]
// 			if err := decoder.Decode(&r); err != nil {
// 				t.Error(err)
// 			}
// 			if !gomock.Eq(r).Matches(test.expectedResponse) {
// 				t.Errorf("response body not match have:%v want:%v", r, test.expectedResponse)
// 			}
// 			bizMock.EXPECT()
// 		})
// 	}
// }

// func TestSampleEntitieHandler_List(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	t.Cleanup(func() {
// 		ctrl.Finish()
// 	})

// 	tests := []struct {
// 		name                string
// 		sampleEntityBizMock func() *mock_sampleusecasev1.MockSampleEntityUsecaseInterface
// 		request             func() *http.Request
// 		expectedResponse    response.Response[[]dto.SampleEntityResponse]
// 		expectedStatusCode  int
// 		ctx                 context.Context
// 	}{
// 		{
// 			name: "success",
// 			sampleEntityBizMock: func() *mock_sampleusecasev1.MockSampleEntityUsecaseInterface {
// 				dsMock := mock_sampleusecasev1.NewMockSampleEntityUsecaseInterface(ctrl)
// 				dsMock.EXPECT().List(gomock.Any()).Return([]*sampleentity.Sample{
// 					{ID: 1, Name: "name1", Text: "text1"}, {ID: 2, Name: "name2", Text: "text2"},
// 				}, nil)
// 				return dsMock
// 			},
// 			request: func() *http.Request {
// 				sampleReq := dto.SampleEntityRequest{
// 					Name: "name",
// 					Text: "text",
// 				}
// 				b, _ := json.Marshal(sampleReq)
// 				r := bytes.NewReader(b)
// 				return httptest.NewRequest(http.MethodPost, "/sample-entities", r)
// 			},
// 			expectedResponse: response.Response[[]dto.SampleEntityResponse]{
// 				Message: "",
// 				Status:  http.StatusOK,
// 				Data: dto.SampleEntityListResponses([]*sampleentity.Sample{
// 					{ID: 1, Name: "name1", Text: "text1"}, {ID: 2, Name: "name2", Text: "text2"},
// 				}),
// 			},
// 			ctx:                context.Background(),
// 			expectedStatusCode: http.StatusOK,
// 		},
// 		{
// 			name: "internal server error",
// 			sampleEntityBizMock: func() *mock_sampleusecasev1.MockSampleEntityUsecaseInterface {
// 				dsMock := mock_sampleusecasev1.NewMockSampleEntityUsecaseInterface(ctrl)
// 				dsMock.EXPECT().List(gomock.Any()).Return(nil, errors.New("error"))
// 				return dsMock
// 			},
// 			request: func() *http.Request {
// 				r := bytes.NewReader([]byte{})
// 				return httptest.NewRequest(http.MethodGet, "/sample-entities", r)
// 			},
// 			expectedResponse: response.Response[[]dto.SampleEntityResponse]{
// 				Message: "internal-error",
// 				Status:  http.StatusInternalServerError,
// 				Data:    nil,
// 			},
// 			ctx:                context.Background(),
// 			expectedStatusCode: http.StatusInternalServerError,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			bizMock := test.sampleEntityBizMock()
// 			handler := NewSampleEntityHandler(slog.New(slog.NewTextHandler(os.Stdout, nil)), bizMock)
// 			recorder := httptest.NewRecorder()
// 			handler.List(recorder, test.request())

// 			if recorder.Code != test.expectedStatusCode {
// 				t.Errorf("status code:%d did not match expected value:%d", recorder.Code, test.expectedStatusCode)
// 			}
// 			decoder := json.NewDecoder(recorder.Body)
// 			r := response.Response[[]dto.SampleEntityResponse]{}
// 			if err := decoder.Decode(&r); err != nil {
// 				t.Error(err)
// 			}
// 			if !gomock.Eq(r).Matches(test.expectedResponse) {
// 				t.Errorf("response body not match have:%v want:%v", r, test.expectedResponse)
// 			}
// 			bizMock.EXPECT()
// 		})
// 	}
// }

// func BenchmarkSampleEntity_List(b *testing.B) {
// 	ctrl := gomock.NewController(b)
// 	seBiz := mock_sampleusecasev1.NewMockSampleEntityUsecaseInterface(ctrl)
// 	seBiz.EXPECT().List(gomock.Any()).Return([]*sampleentity.Sample{
// 		{ID: 1, Name: "name1", Text: "text1"}, {ID: 2, Name: "name2", Text: "text2"},
// 	}, nil).AnyTimes()

// 	sampleReq := dto.SampleEntityRequest{
// 		Name: "name",
// 		Text: "text",
// 	}
// 	bs, _ := json.Marshal(sampleReq)
// 	r := bytes.NewReader(bs)

// 	handler := NewSampleEntityHandler(slog.New(slog.NewTextHandler(os.Stdout, nil)), seBiz)

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		handler.List(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/sample-entities", r))
// 	}
// }

// func TestSampleEntityHandler_Update(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	t.Cleanup(func() {
// 		ctrl.Finish()
// 	})

// 	tests := []struct {
// 		name                string
// 		sampleEntityBizMock func() sampleusecasev1.SampleEntityUsecaseInterface
// 		request             func() *http.Request
// 		expectedStatusCode  int
// 		expectedResponse    response.Response[any]
// 	}{
// 		{
// 			name: "success",
// 			sampleEntityBizMock: func() sampleusecasev1.SampleEntityUsecaseInterface {
// 				dsMock := mock_sampleusecasev1.NewMockSampleEntityUsecaseInterface(ctrl)
// 				dsMock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
// 				return dsMock
// 			},
// 			request: func() *http.Request {
// 				sampleReq := dto.SampleEntityRequest{
// 					Name: "updatedName",
// 					Text: "updatedText",
// 				}
// 				b, _ := json.Marshal(sampleReq)
// 				r := httptest.NewRequest(http.MethodPut, "/entities/1/", bytes.NewReader(b))
// 				r.SetPathValue("id", "1")
// 				return r
// 			},
// 			expectedStatusCode: http.StatusOK,
// 			expectedResponse: response.Response[any]{
// 				Message: "Updated successfully",
// 				Status:  http.StatusOK,
// 				Data:    nil,
// 			},
// 		},
// 		{
// 			name: "invalid request",
// 			sampleEntityBizMock: func() sampleusecasev1.SampleEntityUsecaseInterface {
// 				return mock_sampleusecasev1.NewMockSampleEntityUsecaseInterface(ctrl)
// 			},
// 			request: func() *http.Request {
// 				r := httptest.NewRequest(http.MethodPut, "/entities/abc", nil)
// 				return r
// 			},
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedResponse: response.Response[any]{
// 				Message: "invalid-request",
// 				Status:  http.StatusBadRequest,
// 				Data:    nil,
// 			},
// 		},
// 		{
// 			name: "not found",
// 			sampleEntityBizMock: func() sampleusecasev1.SampleEntityUsecaseInterface {
// 				dsMock := mock_sampleusecasev1.NewMockSampleEntityUsecaseInterface(ctrl)
// 				dsMock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(sampleusecasev1.ErrNotFound)
// 				return dsMock
// 			},
// 			request: func() *http.Request {
// 				sampleReq := dto.SampleEntityRequest{
// 					Name: "updatedName",
// 					Text: "updatedText",
// 				}
// 				b, _ := json.Marshal(sampleReq)
// 				r := httptest.NewRequest(http.MethodPut, "/entities/1/", bytes.NewReader(b))
// 				r.SetPathValue("id", "1")
// 				return r
// 			},
// 			expectedStatusCode: http.StatusNotFound,
// 			expectedResponse: response.Response[any]{
// 				Message: "not-found",
// 				Status:  http.StatusNotFound,
// 				Data:    nil,
// 			},
// 		},
// 		{
// 			name: "internal server error",
// 			sampleEntityBizMock: func() sampleusecasev1.SampleEntityUsecaseInterface {
// 				dsMock := mock_sampleusecasev1.NewMockSampleEntityUsecaseInterface(ctrl)
// 				dsMock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("database error"))
// 				return dsMock
// 			},
// 			request: func() *http.Request {
// 				sampleReq := dto.SampleEntityRequest{
// 					Name: "updatedName",
// 					Text: "updatedText",
// 				}
// 				b, _ := json.Marshal(sampleReq)
// 				r := httptest.NewRequest(http.MethodPut, "/entities/1", bytes.NewReader(b))
// 				r.SetPathValue("id", "1")
// 				return r
// 			},
// 			expectedStatusCode: http.StatusInternalServerError,
// 			expectedResponse: response.Response[any]{
// 				Message: "internal-error",
// 				Status:  http.StatusInternalServerError,
// 				Data:    nil,
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			bizMock := test.sampleEntityBizMock()
// 			handler := NewSampleEntityHandler(slog.New(slog.NewTextHandler(os.Stdout, nil)), bizMock)
// 			recorder := httptest.NewRecorder()
// 			handler.Update(recorder, test.request())

// 			if recorder.Code != test.expectedStatusCode {
// 				t.Errorf("status code:%d did not match expected value:%d", recorder.Code, test.expectedStatusCode)
// 			}

// 			var r response.Response[any]
// 			if err := json.NewDecoder(recorder.Body).Decode(&r); err != nil {
// 				t.Errorf("error decoding response body: %v", err)
// 			}

// 			if !gomock.Eq(r).Matches(test.expectedResponse) {
// 				t.Errorf("response body not match have:%v want:%v", r, test.expectedResponse)
// 			}
// 		})
// 	}
// }

// func BenchmarkSampleEntityHandler_Update(b *testing.B) {
// 	ctrl := gomock.NewController(b)
// 	seBiz := mock_sampleusecasev1.NewMockSampleEntityUsecaseInterface(ctrl)
// 	seBiz.EXPECT().Update(gomock.Any(), gomock.Any()).Times(b.N).Return(nil).AnyTimes()

// 	handler := NewSampleEntityHandler(slog.New(slog.NewTextHandler(os.Stdout, nil)), seBiz)

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		recorder := httptest.NewRecorder()
// 		sampleReq := dto.SampleEntityRequest{
// 			Name: "updatedName",
// 			Text: "updatedText",
// 		}
// 		b, _ := json.Marshal(sampleReq)
// 		req := httptest.NewRequest(http.MethodPut, "/entities/1", bytes.NewReader(b))
// 		handler.Update(recorder, req)
// 	}
// }

// func TestSampleEntityHandler_Delete(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	t.Cleanup(func() {
// 		ctrl.Finish()
// 	})

// 	tests := []struct {
// 		name                string
// 		sampleEntityBizMock func() sampleusecasev1.SampleEntityUsecaseInterface
// 		request             func() *http.Request
// 		expectedStatusCode  int
// 		expectedResponse    response.Response[any]
// 	}{
// 		{
// 			name: "success",
// 			sampleEntityBizMock: func() sampleusecasev1.SampleEntityUsecaseInterface {
// 				dsMock := mock_sampleusecasev1.NewMockSampleEntityUsecaseInterface(ctrl)
// 				dsMock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
// 				return dsMock
// 			},
// 			request: func() *http.Request {
// 				r := httptest.NewRequest(http.MethodDelete, "/entities/1/", bytes.NewReader([]byte{}))
// 				r.SetPathValue("id", "1")
// 				return r
// 			},
// 			expectedStatusCode: http.StatusOK,
// 			expectedResponse: response.Response[any]{
// 				Message: "sample entity deleted",
// 				Status:  http.StatusOK,
// 				Data:    nil,
// 			},
// 		},
// 		{
// 			name: "invalid id",
// 			sampleEntityBizMock: func() sampleusecasev1.SampleEntityUsecaseInterface {
// 				return mock_sampleusecasev1.NewMockSampleEntityUsecaseInterface(ctrl)
// 			},
// 			request: func() *http.Request {
// 				r := httptest.NewRequest(http.MethodDelete, "/entities/abc", nil)
// 				return r
// 			},
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedResponse: response.Response[any]{
// 				Message: "invalid-request",
// 				Status:  http.StatusBadRequest,
// 				Data:    nil,
// 			},
// 		},
// 		{
// 			name: "not found",
// 			sampleEntityBizMock: func() sampleusecasev1.SampleEntityUsecaseInterface {
// 				dsMock := mock_sampleusecasev1.NewMockSampleEntityUsecaseInterface(ctrl)
// 				dsMock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(sampleusecasev1.ErrNotFound)
// 				return dsMock
// 			},
// 			request: func() *http.Request {
// 				r := httptest.NewRequest(http.MethodDelete, "/entities/1/", bytes.NewReader([]byte{}))
// 				r.SetPathValue("id", "1")
// 				return r
// 			},
// 			expectedStatusCode: http.StatusNotFound,
// 			expectedResponse: response.Response[any]{
// 				Message: "not-found",
// 				Status:  http.StatusNotFound,
// 				Data:    nil,
// 			},
// 		},
// 		{
// 			name: "internal server error",
// 			sampleEntityBizMock: func() sampleusecasev1.SampleEntityUsecaseInterface {
// 				dsMock := mock_sampleusecasev1.NewMockSampleEntityUsecaseInterface(ctrl)
// 				dsMock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("database error"))
// 				return dsMock
// 			},
// 			request: func() *http.Request {
// 				r := httptest.NewRequest(http.MethodDelete, "/entities/1/", bytes.NewReader([]byte{}))
// 				r.SetPathValue("id", "1")
// 				return r
// 			},
// 			expectedStatusCode: http.StatusInternalServerError,
// 			expectedResponse: response.Response[any]{
// 				Message: "internal-error",
// 				Status:  http.StatusInternalServerError,
// 				Data:    nil,
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			bizMock := test.sampleEntityBizMock()
// 			handler := NewSampleEntityHandler(slog.New(slog.NewTextHandler(os.Stdout, nil)), bizMock)
// 			recorder := httptest.NewRecorder()
// 			handler.Delete(recorder, test.request())

// 			if recorder.Code != test.expectedStatusCode {
// 				t.Errorf("status code:%d did not match expected value:%d", recorder.Code, test.expectedStatusCode)
// 			}

// 			var r response.Response[any]
// 			if err := json.NewDecoder(recorder.Body).Decode(&r); err != nil {
// 				t.Errorf("error decoding response body: %v", err)
// 			}

// 			if !gomock.Eq(r).Matches(test.expectedResponse) {
// 				t.Errorf("response body not match have:%v want:%v", r, test.expectedResponse)
// 			}
// 		})
// 	}
// }

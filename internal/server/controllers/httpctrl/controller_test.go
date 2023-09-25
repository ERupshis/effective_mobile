package httpctrl

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/helpers"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/server/storage"
	"github.com/erupshis/effective_mobile/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestController_insertPersonHandler(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockBaseStorage(ctrl)
	gomock.InOrder(
		mockStorage.EXPECT().AddPerson(gomock.Any(), gomock.Any()).Return(int64(10), nil),
		mockStorage.EXPECT().AddPerson(gomock.Any(), gomock.Any()).Return(int64(-1), fmt.Errorf("some storage error")),
	)

	ts := httptest.NewServer(Create(mockStorage, log).Route())
	defer ts.Close()

	type fields struct {
		strg storage.BaseStorage
		log  logger.BaseLogger
	}
	type args struct {
		body []byte
	}
	type want struct {
		statusCode int
		body       []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "valid base case",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				body: []byte(`{
					"name": "Aleksandr",
					"surname": "Kozlov",
					"patronymic": "Icanovich",
					"age": 20,
					"gender": "MMMMale",
					"country": "RUSj"
				}`),
			},
			want: want{
				statusCode: http.StatusOK,
				body:       []byte("Added with id: 10"),
			},
		},
		{
			name: "invalid empty request body",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				body: []byte{},
			},
			want: want{
				statusCode: http.StatusUnprocessableEntity,
				body:       []byte{},
			},
		},
		{
			name: "invalid broken json body",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				body: []byte(`{
					"name": "Aleksandr",`,
				),
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       []byte{},
			},
		},
		{
			name: "invalid storage error",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				body: []byte(`{
					"name": "Aleksandr",
					"surname": "Kozlov",
					"patronymic": "Icanovich",
					"age": 20,
					"gender": "MMMMale",
					"country": "RUSj"
				}`),
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       []byte{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewBuffer(tt.args.body)
			req, errReq := http.NewRequest(http.MethodPost, ts.URL, body)
			require.NoError(t, errReq)

			req.Header.Add("Content-Type", "application/json")

			resp, errResp := ts.Client().Do(req)
			require.NoError(t, errResp)
			defer helpers.ExecuteWithLogError(resp.Body.Close, log)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.want.body, respBody)
		})
	}
}

func TestController_deletePersonByIdHandler(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockResult := datastructs.PersonData{
		Name:    "name",
		Age:     0,
		Surname: "surname",
	}

	mockStorage := mocks.NewMockBaseStorage(ctrl)
	gomock.InOrder(
		mockStorage.EXPECT().DeletePersonById(gomock.Any(), gomock.Any()).Return(&mockResult, nil),
		mockStorage.EXPECT().DeletePersonById(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some storage error")),
	)

	ts := httptest.NewServer(Create(mockStorage, log).Route())
	defer ts.Close()

	type fields struct {
		strg storage.BaseStorage
		log  logger.BaseLogger
	}
	type args struct {
		query string
	}
	type want struct {
		statusCode int
		body       []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "valid base case",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				query: "/?id=10",
			},
			want: want{
				statusCode: http.StatusOK,
				body:       []byte("successfully deleted"),
			},
		},
		{
			name: "invalid without id in query",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				query: "/",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body:       []byte{},
			},
		},
		{
			name: "invalid wrong id type",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				query: "/?id=qwe",
			},
			want: want{
				statusCode: http.StatusUnprocessableEntity,
				body:       []byte{},
			},
		},
		{
			name: "invalid storage error",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				query: "/?id=123",
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       []byte{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, errReq := http.NewRequest(http.MethodDelete, ts.URL+tt.args.query, nil)
			require.NoError(t, errReq)

			resp, errResp := ts.Client().Do(req)
			require.NoError(t, errResp)
			defer helpers.ExecuteWithLogError(resp.Body.Close, log)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.want.body, respBody)
		})
	}
}

func TestController_updatePersonByIdHandler(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockResult := datastructs.PersonData{
		Name:    "name",
		Age:     0,
		Surname: "surname",
	}

	mockStorage := mocks.NewMockBaseStorage(ctrl)
	gomock.InOrder(
		mockStorage.EXPECT().UpdatePersonById(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mockResult, nil),
		mockStorage.EXPECT().UpdatePersonById(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some storage error")),
	)

	ts := httptest.NewServer(Create(mockStorage, log).Route())
	defer ts.Close()

	type fields struct {
		strg storage.BaseStorage
		log  logger.BaseLogger
	}
	type args struct {
		query string
		body  []byte
	}
	type want struct {
		statusCode    int
		contentLength string
		contentType   string
		body          []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "valid base case",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				query: "/?id=10",
				body: []byte(`{
					"name": "Aleksandr"
				}`),
			},
			want: want{
				statusCode:    http.StatusOK,
				contentLength: "27",
				contentType:   "text/plain",
				body:          []byte("person with id '10' updated"),
			},
		},
		{
			name: "invalid missing query",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				query: "/",
				body: []byte(`{
					"name": "Aleksandr"
				}`),
			},
			want: want{
				statusCode:    http.StatusBadRequest,
				contentLength: "0",
				contentType:   "",
				body:          []byte{},
			},
		},
		{
			name: "invalid wrong query id type",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				query: "/?id=qwe",
				body: []byte(`{
					"name": "Aleksandr"
				}`),
			},
			want: want{
				statusCode:    http.StatusUnprocessableEntity,
				contentLength: "0",
				contentType:   "",
				body:          []byte{},
			},
		},
		{
			name: "invalid broken json body",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				query: "/?id=123",
				body: []byte(`{
					"nam
				}`),
			},
			want: want{
				statusCode:    http.StatusInternalServerError,
				contentLength: "0",
				contentType:   "",
				body:          []byte{},
			},
		},
		{
			name: "invalid storage error",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				query: "/?id=10",
				body: []byte(`{
					"name": "Aleksandr"
				}`),
			},
			want: want{
				statusCode:    http.StatusInternalServerError,
				contentLength: "0",
				contentType:   "",
				body:          []byte{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewBuffer(tt.args.body)
			req, errReq := http.NewRequest(http.MethodPut, ts.URL+tt.args.query, body)
			require.NoError(t, errReq)

			req.Header.Add("Content-Type", "application/json")

			resp, errResp := ts.Client().Do(req)
			require.NoError(t, errResp)
			defer helpers.ExecuteWithLogError(resp.Body.Close, log)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentLength, resp.Header.Get("Content-Length"))
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.want.body, respBody)
		})
	}
}

func TestController_selectPersonsByFilterHandler(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockResult := []datastructs.PersonData{
		{
			Name:    "name",
			Age:     0,
			Surname: "surname",
		},
	}

	mockStorage := mocks.NewMockBaseStorage(ctrl)
	gomock.InOrder(
		mockStorage.EXPECT().SelectPersons(gomock.Any(), gomock.Any()).Return(mockResult, nil),
		mockStorage.EXPECT().SelectPersons(gomock.Any(), gomock.Any()).Return(mockResult, nil),
		mockStorage.EXPECT().SelectPersons(gomock.Any(), gomock.Any()).Return(mockResult, nil),
		mockStorage.EXPECT().SelectPersons(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some storage error")),
	)

	ts := httptest.NewServer(Create(mockStorage, log).Route())
	defer ts.Close()

	type fields struct {
		strg storage.BaseStorage
		log  logger.BaseLogger
	}
	type args struct {
		query string
	}
	type want struct {
		statusCode    int
		contentLength string
		contentType   string
		body          []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "valid base case",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				query: "/?id=10",
			},
			want: want{
				statusCode:    http.StatusOK,
				contentLength: "61",
				contentType:   "application/json",
				body: []byte(`[
	{
		"id": 0,
		"name": "name",
		"surname": "surname"
	}
]`),
			},
		},
		{
			name: "valid empty query",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				query: "/",
			},
			want: want{
				statusCode:    http.StatusOK,
				contentLength: "61",
				contentType:   "application/json",
				body: []byte(`[
	{
		"id": 0,
		"name": "name",
		"surname": "surname"
	}
]`),
			},
		},
		{
			name: "valid mixed keys",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				query: "/?wrong=123&name=asd",
			},
			want: want{
				statusCode:    http.StatusOK,
				contentLength: "61",
				contentType:   "application/json",
				body: []byte(`[
	{
		"id": 0,
		"name": "name",
		"surname": "surname"
	}
]`),
			},
		},
		{
			name: "invalid wrong query keys only",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				query: "/?wrong=123",
			},
			want: want{
				statusCode:    http.StatusUnprocessableEntity,
				contentLength: "0",
				contentType:   "",
				body:          []byte{},
			},
		},
		{
			name: "invalid storage error",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				query: "/?id=123",
			},
			want: want{
				statusCode:    http.StatusInternalServerError,
				contentLength: "0",
				contentType:   "",
				body:          []byte{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, errReq := http.NewRequest(http.MethodGet, ts.URL+tt.args.query, nil)
			require.NoError(t, errReq)

			resp, errResp := ts.Client().Do(req)
			require.NoError(t, errResp)
			defer helpers.ExecuteWithLogError(resp.Body.Close, log)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentLength, resp.Header.Get("Content-Length"))
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.want.body, respBody)
		})
	}
}

func TestController_badRequestHandler(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockBaseStorage(ctrl)
	gomock.InOrder()

	ts := httptest.NewServer(Create(mockStorage, log).Route())
	defer ts.Close()

	type fields struct {
		strg storage.BaseStorage
		log  logger.BaseLogger
	}
	type args struct {
		query string
	}
	type want struct {
		statusCode    int
		contentLength string
		contentType   string
		body          []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "valid base case",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				query: "/asdf",
			},
			want: want{
				statusCode:    http.StatusBadRequest,
				contentLength: "0",
				contentType:   "",
				body:          []byte{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, errReq := http.NewRequest(http.MethodGet, ts.URL+tt.args.query, nil)
			require.NoError(t, errReq)

			resp, errResp := ts.Client().Do(req)
			require.NoError(t, errResp)
			defer helpers.ExecuteWithLogError(resp.Body.Close, log)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentLength, resp.Header.Get("Content-Length"))
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.want.body, respBody)
		})
	}
}

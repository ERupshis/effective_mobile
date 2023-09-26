package extradatactrl

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/erupshis/effective_mobile/internal/client"
	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
	"github.com/erupshis/effective_mobile/mocks"
	"github.com/golang/mock/gomock"
)

func TestController_RunCancelling(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	type fields struct {
		log logger.BaseLogger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "valid base case",
			fields: fields{
				log: log,
			},
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chIn := make(chan datastructs.ExtraDataFilling, 5)

			chError := make(chan msgbroker.Message, 5)
			chOut := make(chan datastructs.ExtraDataFilling, 5)

			c := &Controller{
				chIn:    chIn,
				chError: chError,
				chOut:   chOut,
				log:     tt.fields.log,
			}

			ctx, cancel := context.WithCancel(tt.args.ctx)
			go c.Run(ctx)

			cancel()
			<-chOut
			<-chError

		})
	}
}

func TestController_RunCloseInputChannel(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	type fields struct {
		log logger.BaseLogger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "valid base case",
			fields: fields{
				log: log,
			},
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chIn := make(chan datastructs.ExtraDataFilling, 5)

			chError := make(chan msgbroker.Message, 5)
			chOut := make(chan datastructs.ExtraDataFilling, 5)

			c := &Controller{
				chIn:    chIn,
				chError: chError,
				chOut:   chOut,
				log:     tt.fields.log,
			}

			go c.Run(tt.args.ctx)

			close(chIn)
			<-chOut
			<-chError

		})
	}
}

func TestController_doRemoteApiRequestCorrectPath(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockResult := []byte(`{"age": 30}`)

	mockBaseClient := mocks.NewMockBaseClient(ctrl)
	gomock.InOrder(
		mockBaseClient.EXPECT().DoGetURIWithQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(http.StatusOK), mockResult, nil),
	)

	type fields struct {
		client client.BaseClient
		log    logger.BaseLogger
	}
	type args struct {
		ctx        context.Context
		serviceUrl string
		dataIn     *datastructs.ExtraDataFilling
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid base case",
			fields: fields{
				client: mockBaseClient,
				log:    log,
			},
			args: args{
				ctx:        context.Background(),
				serviceUrl: "some_url",
				dataIn:     &datastructs.ExtraDataFilling{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chIn := make(chan datastructs.ExtraDataFilling, 5)
			chOut := make(chan datastructs.ExtraDataFilling, 5)
			chError := make(chan msgbroker.Message, 5)

			c := &Controller{
				chIn:    chIn,
				chOut:   chOut,
				chError: chError,
				client:  tt.fields.client,
				log:     tt.fields.log,
			}

			if err := c.doRemoteApiRequest(tt.args.ctx, tt.args.serviceUrl, tt.args.dataIn); (err != nil) != tt.wantErr {
				t.Errorf("doRemoteApiRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestController_doRemoteApiRequestErrorPath(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockResult := []byte(`{"age": 30}`)

	mockBaseClient := mocks.NewMockBaseClient(ctrl)
	gomock.InOrder(
		mockBaseClient.EXPECT().DoGetURIWithQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(http.StatusBadRequest), mockResult, nil),
	)

	type fields struct {
		client client.BaseClient
		log    logger.BaseLogger
	}
	type args struct {
		ctx        context.Context
		serviceUrl string
		dataIn     *datastructs.ExtraDataFilling
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid base case",
			fields: fields{
				client: mockBaseClient,
				log:    log,
			},
			args: args{
				ctx:        context.Background(),
				serviceUrl: "some_url",
				dataIn:     &datastructs.ExtraDataFilling{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chIn := make(chan datastructs.ExtraDataFilling, 5)
			chOut := make(chan datastructs.ExtraDataFilling, 5)
			chError := make(chan msgbroker.Message, 5)

			c := &Controller{
				chIn:    chIn,
				chOut:   chOut,
				chError: chError,
				client:  tt.fields.client,
				log:     tt.fields.log,
			}

			if err := c.doRemoteApiRequest(tt.args.ctx, tt.args.serviceUrl, tt.args.dataIn); (err != nil) != tt.wantErr {
				t.Errorf("doRemoteApiRequest() error = %v, wantErr %v", err, tt.wantErr)
			}

			<-chError
		})
	}
}

func TestController_doRemoteApiRequestInvalid(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockResultError := []byte(`{"error": "some error from api"}`)
	mockResultBrokenJson := []byte(`{"age": 3`)

	mockBaseClient := mocks.NewMockBaseClient(ctrl)
	gomock.InOrder(
		mockBaseClient.EXPECT().DoGetURIWithQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(http.StatusBadRequest), nil, fmt.Errorf("some error")),
		mockBaseClient.EXPECT().DoGetURIWithQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(http.StatusBadRequest), mockResultError, nil),
		mockBaseClient.EXPECT().DoGetURIWithQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(http.StatusOK), mockResultBrokenJson, nil),
	)

	type fields struct {
		client client.BaseClient
		log    logger.BaseLogger
	}
	type args struct {
		ctx        context.Context
		serviceUrl string
		dataIn     *datastructs.ExtraDataFilling
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid base case",
			fields: fields{
				client: mockBaseClient,
				log:    log,
			},
			args: args{
				ctx:        context.Background(),
				serviceUrl: "some_url",
				dataIn:     &datastructs.ExtraDataFilling{},
			},
			wantErr: true,
		},
		{
			name: "invalid error from api",
			fields: fields{
				client: mockBaseClient,
				log:    log,
			},
			args: args{
				ctx:        context.Background(),
				serviceUrl: "some_url",
				dataIn:     &datastructs.ExtraDataFilling{},
			},
			wantErr: true,
		},
		{
			name: "invalid broken json with response 200 from remote api",
			fields: fields{
				client: mockBaseClient,
				log:    log,
			},
			args: args{
				ctx:        context.Background(),
				serviceUrl: "some_url",
				dataIn:     &datastructs.ExtraDataFilling{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chIn := make(chan datastructs.ExtraDataFilling, 5)
			chOut := make(chan datastructs.ExtraDataFilling, 5)
			chError := make(chan msgbroker.Message, 5)

			c := &Controller{
				chIn:    chIn,
				chOut:   chOut,
				chError: chError,
				client:  tt.fields.client,
				log:     tt.fields.log,
			}

			if err := c.doRemoteApiRequest(tt.args.ctx, tt.args.serviceUrl, tt.args.dataIn); (err != nil) != tt.wantErr {
				t.Errorf("doRemoteApiRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

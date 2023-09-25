package msgsavectrl

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
	"github.com/erupshis/effective_mobile/internal/server/storage"
	"github.com/erupshis/effective_mobile/mocks"
	"github.com/golang/mock/gomock"
)

func TestController_Run(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockBaseStorage(ctrl)
	gomock.InOrder(
		mockStorage.EXPECT().AddPerson(gomock.Any(), gomock.Any()).Return(int64(72), nil),
	)

	type fields struct {
		strg storage.BaseStorage
		log  logger.BaseLogger
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
			name: "valid base case. should log successful save.",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		chIn := make(chan datastructs.ExtraDataFilling, 5)
		chError := make(chan msgbroker.Message, 5)

		c := &Controller{
			chIn:    chIn,
			chError: chError,
			strg:    tt.fields.strg,
			log:     tt.fields.log,
		}

		chIn <- datastructs.ExtraDataFilling{}

		go c.Run(tt.args.ctx)

		time.Sleep(time.Second)
	}
}

func TestController_RunFail(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockBaseStorage(ctrl)
	gomock.InOrder(
		mockStorage.EXPECT().AddPerson(gomock.Any(), gomock.Any()).Return(int64(-1), fmt.Errorf("wrong save")),
	)

	type fields struct {
		strg storage.BaseStorage
		log  logger.BaseLogger
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
			name: "valid base case. should log successful save.",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		chIn := make(chan datastructs.ExtraDataFilling, 5)
		chError := make(chan msgbroker.Message, 5)

		c := &Controller{
			chIn:    chIn,
			chError: chError,
			strg:    tt.fields.strg,
			log:     tt.fields.log,
		}

		chIn <- datastructs.ExtraDataFilling{}

		go c.Run(tt.args.ctx)
		<-chError
	}
}

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

			c := &Controller{
				chIn:    chIn,
				chError: chError,
				log:     tt.fields.log,
			}

			ctx, cancel := context.WithCancel(tt.args.ctx)
			go c.Run(ctx)

			cancel()
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

			c := &Controller{
				chIn:    chIn,
				chError: chError,
				log:     tt.fields.log,
			}

			go c.Run(tt.args.ctx)

			close(chIn)
			<-chError
		})
	}
}

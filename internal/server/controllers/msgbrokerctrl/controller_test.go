package msgbrokerctrl

import (
	"context"
	"reflect"
	"testing"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
	"github.com/stretchr/testify/assert"
)

func TestController_handleMessage(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	type fields struct {
		log logger.BaseLogger
	}
	type args struct {
		msg msgbroker.Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *datastructs.PersonData
		wantErr bool
	}{
		{
			name: "valid base case",
			fields: fields{
				log: log,
			},
			args: args{
				msg: msgbroker.Message{
					Key:   []byte("person data"),
					Value: []byte(`{"name":"Ekaterina","surname":"Ivanova","patronymic":"Sergeevna"}`),
				},
			},
			want: &datastructs.PersonData{
				Name:       "Ekaterina",
				Surname:    "Ivanova",
				Patronymic: "Sergeevna",
			},
			wantErr: false,
		},
		{
			name: "valid base case without patronymic",
			fields: fields{
				log: log,
			},
			args: args{
				msg: msgbroker.Message{
					Key:   []byte("person data"),
					Value: []byte(`{"name":"Ekaterina","surname":"Ivanova"}`),
				},
			},
			want: &datastructs.PersonData{
				Name:    "Ekaterina",
				Surname: "Ivanova",
			},
			wantErr: false,
		},
		{
			name: "invalid without critical field",
			fields: fields{
				log: log,
			},
			args: args{
				msg: msgbroker.Message{
					Key:   []byte("person data"),
					Value: []byte(`{"name":"Ekaterina"}`),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid broken json",
			fields: fields{
				log: log,
			},
			args: args{
				msg: msgbroker.Message{
					Key:   []byte("person data"),
					Value: []byte(`{"name":"Ekaterina","surname":"Ivanova"`),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chIn := make(chan msgbroker.Message, 5)
			chError := make(chan msgbroker.Message, 5)
			chOut := make(chan datastructs.ExtraDataFilling, 5)

			c := &Controller{
				chIn:    chIn,
				chError: chError,
				chOut:   chOut,
				log:     tt.fields.log,
			}
			got, err := c.handleMessage(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestController_Run(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	type fields struct {
		log logger.BaseLogger
	}
	type args struct {
		ctx context.Context
		msg msgbroker.Message
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantMsgOk  bool
		wantMsgErr bool
	}{
		{
			name: "valid base case",
			fields: fields{
				log: log,
			},
			args: args{
				ctx: context.Background(),
				msg: msgbroker.Message{
					Key:   []byte("person data"),
					Value: []byte(`{"name":"Ekaterina","surname":"Ivanova","patronymic":"Sergeevna"}`),
				},
			},
			wantMsgOk:  true,
			wantMsgErr: false,
		},
		{
			name: "invalid message without critical field",
			fields: fields{
				log: log,
			},
			args: args{
				ctx: context.Background(),
				msg: msgbroker.Message{
					Key:   []byte("person data"),
					Value: []byte(`{"name":"Ekaterina","patronymic":"Sergeevna"}`),
				},
			},
			wantMsgOk:  false,
			wantMsgErr: true,
		},
		{
			name: "invalid message with broken json",
			fields: fields{
				log: log,
			},
			args: args{
				ctx: context.Background(),
				msg: msgbroker.Message{
					Key:   []byte("person data"),
					Value: []byte(`{"name":"Ekaterina","surname":"Ivanova""patronymic":"Sergeevna"}`),
				},
			},
			wantMsgOk:  false,
			wantMsgErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chIn := make(chan msgbroker.Message, 5)
			chIn <- tt.args.msg

			chError := make(chan msgbroker.Message, 5)
			chOut := make(chan datastructs.ExtraDataFilling, 5)

			c := &Controller{
				chIn:    chIn,
				chError: chError,
				chOut:   chOut,
				log:     tt.fields.log,
			}

			ctx, cancel := context.WithCancel(tt.args.ctx)
			defer cancel()
			go c.Run(ctx)

			if tt.wantMsgOk {
				_, ok := <-chOut
				assert.True(t, ok, "check channel is open")
			}

			if tt.wantMsgErr {
				_, ok := <-chError
				assert.True(t, ok, "check channel is open")
			}
		})
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
			chIn := make(chan msgbroker.Message, 5)

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
			chIn := make(chan msgbroker.Message, 5)

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

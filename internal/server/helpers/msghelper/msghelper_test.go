package msghelper

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
)

func TestCreateErrorMessage(t *testing.T) {
	type args struct {
		originalMsg []byte
		err         error
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "valid base case",
			args: args{
				originalMsg: []byte("some-error"),
				err:         fmt.Errorf("some error text"),
			},
			want:    []byte(`{"error":"some error text","original":"some-error"}`),
			wantErr: false,
		},
		{
			name: "valid with empty original message",
			args: args{
				originalMsg: nil,
				err:         fmt.Errorf("some error text"),
			},
			want:    []byte(`{"error":"some error text","original":""}`),
			wantErr: false,
		},
		{
			name: "invalid nil error",
			args: args{
				originalMsg: []byte("some-error"),
				err:         nil,
			},
			want:    []byte{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateErrorMessage(tt.args.originalMsg, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateErrorMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateErrorMessage() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestIsMessageValid(t *testing.T) {
	type args struct {
		personData datastructs.PersonData
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid base case",
			args: args{
				personData: datastructs.PersonData{
					Name:       "asd",
					Surname:    "asd",
					Patronymic: "asd",
					Age:        30,
					Gender:     "asd",
					Country:    "asd",
				}},
			want:    true,
			wantErr: false,
		},
		{
			name: "valid without patronymic",
			args: args{
				personData: datastructs.PersonData{
					Name:    "asd",
					Surname: "asd",
					Age:     30,
					Gender:  "asd",
					Country: "asd",
				}},
			want:    true,
			wantErr: false,
		},
		{
			name: "valid without non-critical field",
			args: args{
				personData: datastructs.PersonData{
					Name:       "asd",
					Surname:    "asd",
					Patronymic: "asd",
					Age:        30,
					Gender:     "asd",
				}},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid without critical field name",
			args: args{
				personData: datastructs.PersonData{
					Surname:    "asd",
					Patronymic: "asd",
					Age:        30,
					Gender:     "asd",
					Country:    "asd",
				}},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid without critical field surname",
			args: args{
				personData: datastructs.PersonData{
					Name:       "asd",
					Patronymic: "asd",
					Age:        30,
					Gender:     "asd",
					Country:    "asd",
				}},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid without critical fields",
			args: args{
				personData: datastructs.PersonData{
					Patronymic: "asd",
					Age:        30,
					Gender:     "asd",
					Country:    "asd",
				}},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsMessageValid(tt.args.personData)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsMessageValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsMessageValid() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPutErrorMessageInChan(t *testing.T) {
	type args struct {
		chError   chan msgbroker.Message
		msg       *msgbroker.Message
		errMsgKey string
		err       error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid base case",
			args: args{
				chError: make(chan msgbroker.Message),
				msg: &msgbroker.Message{
					Key:   []byte{},
					Value: []byte{},
				},
				errMsgKey: "msgKey",
				err:       fmt.Errorf("some error text"),
			},
			wantErr: false,
		},
		{
			name: "invalid nil channel",
			args: args{
				chError: nil,
				msg: &msgbroker.Message{
					Key:   []byte{},
					Value: []byte{},
				},
				errMsgKey: "msgKey",
				err:       fmt.Errorf("some error text"),
			},
			wantErr: true,
		},
		{
			name: "invalid nil message",
			args: args{
				chError:   make(chan msgbroker.Message),
				msg:       nil,
				errMsgKey: "msgKey",
				err:       fmt.Errorf("some error text"),
			},
			wantErr: true,
		},
		{
			name: "invalid nil channel and message",
			args: args{
				chError:   nil,
				msg:       nil,
				errMsgKey: "msgKey",
				err:       fmt.Errorf("some error text"),
			},
			wantErr: true,
		},
		{
			name: "invalid nil err",
			args: args{
				chError: make(chan msgbroker.Message),
				msg: &msgbroker.Message{
					Key:   []byte{},
					Value: []byte{},
				},
				errMsgKey: "msgKey",
				err:       nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.chError != nil {
				go func() {
					select {
					case <-tt.args.chError:
						return
					}
				}()
			}

			if err := PutErrorMessageInChan(tt.args.chError, tt.args.msg, tt.args.errMsgKey, tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("PutErrorMessageInChan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

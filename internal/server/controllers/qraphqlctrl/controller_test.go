package qraphqlctrl

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/server/storage"
	"github.com/erupshis/effective_mobile/mocks"
	"github.com/golang/mock/gomock"
	"github.com/graphql-go/graphql"
)

func TestController_selectPersonsResolver(t *testing.T) {
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
		mockStorage.EXPECT().SelectPersons(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some storage error")),
	)

	type fields struct {
		strg storage.BaseStorage
		log  logger.BaseLogger
	}
	type args struct {
		p graphql.ResolveParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "valid base case",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				p: graphql.ResolveParams{},
			},
			want: []datastructs.PersonData{
				{
					Name:    "name",
					Age:     0,
					Surname: "surname",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid result - some incorrect input or problem on storage",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				p: graphql.ResolveParams{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{
				strg: tt.fields.strg,
				log:  tt.fields.log,
			}
			got, err := c.selectPersonsResolver(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("selectPersonsResolver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("selectPersonsResolver() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestController_insertPersonResolver(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockBaseStorage(ctrl)
	gomock.InOrder(
		mockStorage.EXPECT().AddPerson(gomock.Any(), gomock.Any()).Return(int64(10), nil),
		mockStorage.EXPECT().AddPerson(gomock.Any(), gomock.Any()).Return(int64(-1), fmt.Errorf("some storage error")),
	)

	type fields struct {
		strg storage.BaseStorage
		log  logger.BaseLogger
	}
	type args struct {
		p graphql.ResolveParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "valid base case",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				p: graphql.ResolveParams{},
			},
			want: datastructs.PersonData{
				Id: 10,
			},
			wantErr: false,
		},
		{
			name: "invalid result - some incorrect input or problem on storage",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				p: graphql.ResolveParams{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{
				strg: tt.fields.strg,
				log:  tt.fields.log,
			}
			got, err := c.insertPersonResolver(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("insertPersonResolver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("insertPersonResolver() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestController_deletePersonResolver(t *testing.T) {
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

	type fields struct {
		strg storage.BaseStorage
		log  logger.BaseLogger
	}
	type args struct {
		p graphql.ResolveParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "valid base case",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				p: graphql.ResolveParams{},
			},
			want: &datastructs.PersonData{
				Name:    "name",
				Age:     0,
				Surname: "surname",
			},
			wantErr: false,
		},
		{
			name: "invalid result - some incorrect input or problem on storage",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				p: graphql.ResolveParams{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{
				strg: tt.fields.strg,
				log:  tt.fields.log,
			}
			got, err := c.deletePersonResolver(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("updatePersonResolver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updatePersonResolver() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestController_updatePersonResolver(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockResult := datastructs.PersonData{
		Id:      10,
		Name:    "name",
		Age:     0,
		Surname: "surname",
	}

	mockStorage := mocks.NewMockBaseStorage(ctrl)
	gomock.InOrder(
		mockStorage.EXPECT().UpdatePersonById(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mockResult, nil),
		mockStorage.EXPECT().UpdatePersonById(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some storage error")),
	)

	type fields struct {
		strg storage.BaseStorage
		log  logger.BaseLogger
	}
	type args struct {
		p graphql.ResolveParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "valid base case",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				p: graphql.ResolveParams{},
			},
			want: &datastructs.PersonData{
				Id:      10,
				Name:    "name",
				Age:     0,
				Surname: "surname",
			},
			wantErr: false,
		},
		{
			name: "invalid result - some incorrect input or problem on storage",
			fields: fields{
				strg: mockStorage,
				log:  log,
			},
			args: args{
				p: graphql.ResolveParams{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{
				strg: tt.fields.strg,
				log:  tt.fields.log,
			}
			got, err := c.updatePersonResolver(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("updatePersonResolver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updatePersonResolver() got = %v, want %v", got, tt.want)
			}
		})
	}
}

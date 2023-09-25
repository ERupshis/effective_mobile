package storage

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/server/storage/managers"
	"github.com/erupshis/effective_mobile/mocks"
	"github.com/golang/mock/gomock"
)

func Test_storage_AddPerson(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBaseStorageManager := mocks.NewMockBaseStorageManager(ctrl)
	gomock.InOrder(
		mockBaseStorageManager.EXPECT().AddPerson(gomock.Any(), gomock.Any()).Return(int64(10), nil),
		mockBaseStorageManager.EXPECT().AddPerson(gomock.Any(), gomock.Any()).Return(int64(-1), fmt.Errorf("some storage error")),
	)

	type fields struct {
		manager managers.BaseStorageManager
		log     logger.BaseLogger
	}
	type args struct {
		ctx       context.Context
		newPerson *datastructs.PersonData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "valid base case",
			fields: fields{
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx: context.Background(),
				newPerson: &datastructs.PersonData{
					Name:    "name",
					Surname: "suranme",
					Age:     12,
					Gender:  "male",
					Country: "ru",
				},
			},
			want:    10,
			wantErr: false,
		},
		{
			name: "invalid not full person's data",
			fields: fields{
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx: context.Background(),
				newPerson: &datastructs.PersonData{
					Name:    "name",
					Surname: "suranme",
					Age:     12,
				},
			},
			want:    -1,
			wantErr: true,
		},
		{
			name: "invalid storage manager error",
			fields: fields{
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx: context.Background(),
				newPerson: &datastructs.PersonData{
					Name:    "name",
					Surname: "suranme",
					Age:     12,
					Gender:  "male",
					Country: "ru",
				},
			},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				manager: tt.fields.manager,
				log:     tt.fields.log,
			}
			got, err := s.AddPerson(tt.args.ctx, tt.args.newPerson)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddPerson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddPerson() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_storage_SelectPersons(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockResult := []datastructs.PersonData{
		{
			Name: "name",
		},
	}

	mockBaseStorageManager := mocks.NewMockBaseStorageManager(ctrl)
	gomock.InOrder(
		mockBaseStorageManager.EXPECT().SelectPersons(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResult, nil),
		mockBaseStorageManager.EXPECT().SelectPersons(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResult, nil),
		mockBaseStorageManager.EXPECT().SelectPersons(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]datastructs.PersonData{}, nil),
		mockBaseStorageManager.EXPECT().SelectPersons(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some storage error")),
		mockBaseStorageManager.EXPECT().SelectPersons(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some storage error")),
	)

	type fields struct {
		manager managers.BaseStorageManager
		log     logger.BaseLogger
	}
	type args struct {
		ctx    context.Context
		values map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []datastructs.PersonData
		wantErr bool
	}{
		{
			name: "valid base case",
			fields: fields{
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx:    context.Background(),
				values: map[string]interface{}{},
			},
			want: []datastructs.PersonData{
				{
					Name: "name",
				},
			},
			wantErr: false,
		},
		{
			name: "valid query has id",
			fields: fields{
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx: context.Background(),
				values: map[string]interface{}{
					"id": "12",
				},
			},
			want: []datastructs.PersonData{
				{
					Name: "name",
				},
			},
			wantErr: false,
		},
		{
			name: "valid no one person's data didn't satisfies filters",
			fields: fields{
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx: context.Background(),
				values: map[string]interface{}{
					"id": "12",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid case no matches to filter",
			fields: fields{
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx:    context.Background(),
				values: map[string]interface{}{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid negative page_num",
			fields: fields{
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx: context.Background(),
				values: map[string]interface{}{
					"page_num": "-1",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid negative page_size",
			fields: fields{
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx: context.Background(),
				values: map[string]interface{}{
					"page_size": "-1",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid storage manager error",
			fields: fields{
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx:    context.Background(),
				values: map[string]interface{}{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				manager: tt.fields.manager,
				log:     tt.fields.log,
			}
			got, err := s.SelectPersons(tt.args.ctx, tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectPersons() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelectPersons() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_storage_DeletePersonById(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockResult := []datastructs.PersonData{
		{
			Name: "name",
		},
	}

	mockBaseStorageManager := mocks.NewMockBaseStorageManager(ctrl)
	gomock.InOrder(
		mockBaseStorageManager.EXPECT().SelectPersons(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResult, nil),
		mockBaseStorageManager.EXPECT().DeletePersonById(gomock.Any(), gomock.Any()).Return(int64(10), nil),
		mockBaseStorageManager.EXPECT().SelectPersons(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some storage error")),
		mockBaseStorageManager.EXPECT().SelectPersons(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResult, nil),
		mockBaseStorageManager.EXPECT().DeletePersonById(gomock.Any(), gomock.Any()).Return(int64(-1), fmt.Errorf("some storage error")),
		mockBaseStorageManager.EXPECT().SelectPersons(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResult, nil),
		mockBaseStorageManager.EXPECT().DeletePersonById(gomock.Any(), gomock.Any()).Return(int64(0), nil),
	)

	type fields struct {
		manager managers.BaseStorageManager
		log     logger.BaseLogger
	}
	type args struct {
		ctx context.Context
		id  int64
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
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx: context.Background(),
				id:  10,
			},
			want: &datastructs.PersonData{
				Name: "name",
			},
			wantErr: false,
		},
		{
			name: "invalid case storage manager error on selectById",
			fields: fields{
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx: context.Background(),
				id:  10,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid case storage manager error on deleteById",
			fields: fields{
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx: context.Background(),
				id:  10,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid case storage manager didn't remove any person data by id",
			fields: fields{
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx: context.Background(),
				id:  10,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				manager: tt.fields.manager,
				log:     tt.fields.log,
			}
			got, err := s.DeletePersonById(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeletePersonById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeletePersonById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_storage_UpdatePersonById(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockResult := []datastructs.PersonData{
		{
			Name: "name",
		},
	}

	mockBaseStorageManager := mocks.NewMockBaseStorageManager(ctrl)
	gomock.InOrder(
		mockBaseStorageManager.EXPECT().UpdatePersonById(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(10), nil),
		mockBaseStorageManager.EXPECT().SelectPersons(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResult, nil),
		mockBaseStorageManager.EXPECT().UpdatePersonById(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(10), nil),
		mockBaseStorageManager.EXPECT().SelectPersons(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResult, nil),
		mockBaseStorageManager.EXPECT().UpdatePersonById(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(10), fmt.Errorf("some error")),
		mockBaseStorageManager.EXPECT().UpdatePersonById(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), nil),
		mockBaseStorageManager.EXPECT().UpdatePersonById(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(10), nil),
		mockBaseStorageManager.EXPECT().SelectPersons(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResult, fmt.Errorf("some error")),
	)

	type fields struct {
		manager managers.BaseStorageManager
		log     logger.BaseLogger
	}
	type args struct {
		ctx    context.Context
		id     int64
		values map[string]interface{}
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
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx: context.Background(),
				id:  10,
				values: map[string]interface{}{
					"name": "new_name",
				},
			},
			want: &datastructs.PersonData{
				Name: "name",
			},
			wantErr: false,
		},
		{
			name: "valid mix of filters",
			fields: fields{
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx: context.Background(),
				id:  10,
				values: map[string]interface{}{
					"name":         "new_name",
					"surnameWrong": "surname",
				},
			},
			want: &datastructs.PersonData{
				Name: "name",
			},
			wantErr: false,
		},
		{
			name: "invalid missing correct filters",
			fields: fields{
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx: context.Background(),
				id:  10,
				values: map[string]interface{}{
					"nameWrong": "new_name",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid storage manager error on Update",
			fields: fields{
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx: context.Background(),
				id:  10,
				values: map[string]interface{}{
					"name": "new_name",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid storage manager didn't update any person data",
			fields: fields{
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx: context.Background(),
				id:  10,
				values: map[string]interface{}{
					"name": "new_name",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid storage maanger couldn't select modified person's data by id",
			fields: fields{
				manager: mockBaseStorageManager,
				log:     log,
			},
			args: args{
				ctx: context.Background(),
				id:  10,
				values: map[string]interface{}{
					"name": "new_name",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				manager: tt.fields.manager,
				log:     tt.fields.log,
			}
			got, err := s.UpdatePersonById(tt.args.ctx, tt.args.id, tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePersonById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdatePersonById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

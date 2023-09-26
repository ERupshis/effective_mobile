package cache

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/server/storage"
	"github.com/erupshis/effective_mobile/internal/server/storage/cache/manager"
	"github.com/erupshis/effective_mobile/mocks"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
)

func Test_cache_AddPerson(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBaseCacheManager := mocks.NewMockBaseCacheManager(ctrl)
	gomock.InOrder(
		mockBaseCacheManager.EXPECT().Flush(gomock.Any()).Return(nil),
		mockBaseCacheManager.EXPECT().Flush(gomock.Any()).Return(fmt.Errorf("some error")),
	)

	mockBaseStorage := mocks.NewMockBaseStorage(ctrl)
	gomock.InOrder(
		mockBaseStorage.EXPECT().AddPerson(gomock.Any(), gomock.Any()).Return(int64(10), nil),
		mockBaseStorage.EXPECT().AddPerson(gomock.Any(), gomock.Any()).Return(int64(10), nil),
		mockBaseStorage.EXPECT().AddPerson(gomock.Any(), gomock.Any()).Return(int64(-1), fmt.Errorf("some error")),
	)

	type fields struct {
		manager manager.BaseCacheManager
		storage storage.BaseStorage
		logger  logger.BaseLogger
	}
	type args struct {
		ctx  context.Context
		data *datastructs.PersonData
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
				manager: mockBaseCacheManager,
				storage: mockBaseStorage,
				logger:  log,
			},
			args: args{
				ctx: context.Background(),
				data: &datastructs.PersonData{
					Name: "name",
				},
			},
			want:    10,
			wantErr: false,
		},
		{
			name: "valid with error on flushing",
			fields: fields{
				manager: mockBaseCacheManager,
				storage: mockBaseStorage,
				logger:  log,
			},
			args: args{
				ctx: context.Background(),
				data: &datastructs.PersonData{
					Name: "name",
				},
			},
			want:    10,
			wantErr: true,
		},
		{
			name: "invalid with storage error",
			fields: fields{
				manager: mockBaseCacheManager,
				storage: mockBaseStorage,
				logger:  log,
			},
			args: args{
				ctx: context.Background(),
				data: &datastructs.PersonData{
					Name: "name",
				},
			},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				manager: tt.fields.manager,
				storage: tt.fields.storage,
				logger:  tt.fields.logger,
			}
			got, err := c.AddPerson(tt.args.ctx, tt.args.data)
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

func Test_cache_SelectPersons(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockResultRaw := []byte(`[{"id":77,"name":"Olga","country":"RUS"}]`)
	mockResultRawBroken := []byte(`[{"id":77,"name":"Olga","country":"RU}]`)
	mockResultData := []datastructs.PersonData{
		{
			Id:      77,
			Name:    "Olga",
			Country: "RUS",
		},
	}

	mockBaseCacheManager := mocks.NewMockBaseCacheManager(ctrl)
	gomock.InOrder(
		//"valid base case with value in cache"
		mockBaseCacheManager.EXPECT().Get(gomock.Any(), gomock.Any()).Return(mockResultRaw, nil),
		//"valid base case with value in cache but raw json is damaged"
		mockBaseCacheManager.EXPECT().Get(gomock.Any(), gomock.Any()).Return(mockResultRawBroken, nil),
		mockBaseCacheManager.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		//"valid case, cache manager returns Nil"
		mockBaseCacheManager.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, redis.Nil),
		mockBaseCacheManager.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		//"valid case, cache manager returns unknown error"
		mockBaseCacheManager.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some error")),
		mockBaseCacheManager.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		//"invalid storage error"
		mockBaseCacheManager.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some error")),
		//"invalid cache manager error on add"
		mockBaseCacheManager.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some error")),
		//mockBaseCacheManager.EXPECT().Flush(gomock.Any()).Return(fmt.Errorf("some error")),
	)

	mockBaseStorage := mocks.NewMockBaseStorage(ctrl)
	gomock.InOrder(
		//"valid base case with value in cache"
		mockBaseStorage.EXPECT().SelectPersons(gomock.Any(), gomock.Any()).Return(mockResultData, nil),
		//"valid base case with value in cache but raw json is damaged"
		mockBaseStorage.EXPECT().SelectPersons(gomock.Any(), gomock.Any()).Return(mockResultData, nil),
		//"valid case, cache manager returns Nil",
		mockBaseStorage.EXPECT().SelectPersons(gomock.Any(), gomock.Any()).Return(mockResultData, nil),
		//"invalid storage error"
		mockBaseStorage.EXPECT().SelectPersons(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some error")),
		//"invalid cache manager error on add"
		mockBaseStorage.EXPECT().SelectPersons(gomock.Any(), gomock.Any()).Return(mockResultData, nil),
		mockBaseCacheManager.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error")),
	)

	type fields struct {
		manager manager.BaseCacheManager
		storage storage.BaseStorage
		logger  logger.BaseLogger
	}
	type args struct {
		ctx     context.Context
		filters map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []datastructs.PersonData
		wantErr bool
	}{
		{
			name: "valid base case with value in cache",
			fields: fields{
				manager: mockBaseCacheManager,
				storage: mockBaseStorage,
				logger:  log,
			},
			args: args{
				ctx: context.Background(),
				filters: map[string]interface{}{
					"name": "asd",
				},
			},
			want: []datastructs.PersonData{
				{
					Id:      77,
					Name:    "Olga",
					Country: "RUS",
				},
			},
			wantErr: false,
		},
		{
			name: "valid base case with value in cache but raw json is damaged",
			fields: fields{
				manager: mockBaseCacheManager,
				storage: mockBaseStorage,
				logger:  log,
			},
			args: args{
				ctx: context.Background(),
				filters: map[string]interface{}{
					"name": "asd",
				},
			},
			want: []datastructs.PersonData{
				{
					Id:      77,
					Name:    "Olga",
					Country: "RUS",
				},
			},
			wantErr: false,
		},
		{
			name: "valid case, cache manager returns Nil",
			fields: fields{
				manager: mockBaseCacheManager,
				storage: mockBaseStorage,
				logger:  log,
			},
			args: args{
				ctx: context.Background(),
				filters: map[string]interface{}{
					"name": "asd",
				},
			},
			want: []datastructs.PersonData{
				{
					Id:      77,
					Name:    "Olga",
					Country: "RUS",
				},
			},
			wantErr: false,
		},
		{
			name: "valid case, cache manager returns unknown error",
			fields: fields{
				manager: mockBaseCacheManager,
				storage: mockBaseStorage,
				logger:  log,
			},
			args: args{
				ctx: context.Background(),
				filters: map[string]interface{}{
					"name": "asd",
				},
			},
			want: []datastructs.PersonData{
				{
					Id:      77,
					Name:    "Olga",
					Country: "RUS",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid storage error",
			fields: fields{
				manager: mockBaseCacheManager,
				storage: mockBaseStorage,
				logger:  log,
			},
			args: args{
				ctx: context.Background(),
				filters: map[string]interface{}{
					"name": "asd",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid cache manager error on add",
			fields: fields{
				manager: mockBaseCacheManager,
				storage: mockBaseStorage,
				logger:  log,
			},
			args: args{
				ctx: context.Background(),
				filters: map[string]interface{}{
					"name": "asd",
				},
			},
			want: []datastructs.PersonData{
				{
					Id:      77,
					Name:    "Olga",
					Country: "RUS",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				manager: tt.fields.manager,
				storage: tt.fields.storage,
				logger:  tt.fields.logger,
			}
			got, err := c.SelectPersons(tt.args.ctx, tt.args.filters)
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

func Test_cache_DeletePersonById(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBaseCacheManager := mocks.NewMockBaseCacheManager(ctrl)
	gomock.InOrder(
		mockBaseCacheManager.EXPECT().Flush(gomock.Any()).Return(nil),
		mockBaseCacheManager.EXPECT().Flush(gomock.Any()).Return(fmt.Errorf("some error")),
	)

	mockResultData := &datastructs.PersonData{
		Id:      77,
		Name:    "Olga",
		Country: "RUS",
	}

	mockBaseStorage := mocks.NewMockBaseStorage(ctrl)
	gomock.InOrder(
		mockBaseStorage.EXPECT().DeletePersonById(gomock.Any(), gomock.Any()).Return(mockResultData, nil),
		mockBaseStorage.EXPECT().DeletePersonById(gomock.Any(), gomock.Any()).Return(mockResultData, nil),
		mockBaseStorage.EXPECT().DeletePersonById(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some error")),
	)

	type fields struct {
		manager manager.BaseCacheManager
		storage storage.BaseStorage
		logger  logger.BaseLogger
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
				manager: mockBaseCacheManager,
				storage: mockBaseStorage,
				logger:  log,
			},
			args: args{
				ctx: context.Background(),
				id:  77,
			},
			want: &datastructs.PersonData{
				Id:      77,
				Name:    "Olga",
				Country: "RUS",
			},
			wantErr: false,
		},
		{
			name: "valid with error on flushing",
			fields: fields{
				manager: mockBaseCacheManager,
				storage: mockBaseStorage,
				logger:  log,
			},
			args: args{
				ctx: context.Background(),
				id:  77,
			},
			want: &datastructs.PersonData{
				Id:      77,
				Name:    "Olga",
				Country: "RUS",
			},
			wantErr: true,
		},
		{
			name: "invalid with storage error",
			fields: fields{
				manager: mockBaseCacheManager,
				storage: mockBaseStorage,
				logger:  log,
			},
			args: args{
				ctx: context.Background(),
				id:  77,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				manager: tt.fields.manager,
				storage: tt.fields.storage,
				logger:  tt.fields.logger,
			}
			got, err := c.DeletePersonById(tt.args.ctx, tt.args.id)
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

func Test_cache_UpdatePersonById(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBaseCacheManager := mocks.NewMockBaseCacheManager(ctrl)
	gomock.InOrder(
		mockBaseCacheManager.EXPECT().Flush(gomock.Any()).Return(nil),
		mockBaseCacheManager.EXPECT().Flush(gomock.Any()).Return(fmt.Errorf("some error")),
	)

	mockResultData := &datastructs.PersonData{
		Id:      77,
		Name:    "Olga",
		Country: "RUS",
	}

	mockBaseStorage := mocks.NewMockBaseStorage(ctrl)
	gomock.InOrder(
		mockBaseStorage.EXPECT().UpdatePersonById(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResultData, nil),
		mockBaseStorage.EXPECT().UpdatePersonById(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResultData, nil),
		mockBaseStorage.EXPECT().UpdatePersonById(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some error")),
	)

	type fields struct {
		manager manager.BaseCacheManager
		storage storage.BaseStorage
		logger  logger.BaseLogger
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
				manager: mockBaseCacheManager,
				storage: mockBaseStorage,
				logger:  log,
			},
			args: args{
				ctx: context.Background(),
				id:  77,
			},
			want: &datastructs.PersonData{
				Id:      77,
				Name:    "Olga",
				Country: "RUS",
			},
			wantErr: false,
		},
		{
			name: "valid with error on flushing",
			fields: fields{
				manager: mockBaseCacheManager,
				storage: mockBaseStorage,
				logger:  log,
			},
			args: args{
				ctx: context.Background(),
				id:  77,
			},
			want: &datastructs.PersonData{
				Id:      77,
				Name:    "Olga",
				Country: "RUS",
			},
			wantErr: true,
		},
		{
			name: "invalid with storage error",
			fields: fields{
				manager: mockBaseCacheManager,
				storage: mockBaseStorage,
				logger:  log,
			},
			args: args{
				ctx: context.Background(),
				id:  77,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				manager: tt.fields.manager,
				storage: tt.fields.storage,
				logger:  tt.fields.logger,
			}
			got, err := c.UpdatePersonById(tt.args.ctx, tt.args.id, tt.args.values)
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

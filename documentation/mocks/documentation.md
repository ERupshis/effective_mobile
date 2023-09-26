package mocks // import "github.com/erupshis/effective_mobile/mocks"

Package mocks is a generated GoMock package.

Package mocks is a generated GoMock package.

Package mocks is a generated GoMock package.

Package mocks is a generated GoMock package.

Package mocks is a generated GoMock package.

Package mocks is a generated GoMock package.

TYPES

type MockBaseCacheManager struct {
	// Has unexported fields.
}
    MockBaseCacheManager is a mock of BaseCacheManager interface.

func NewMockBaseCacheManager(ctrl *gomock.Controller) *MockBaseCacheManager
    NewMockBaseCacheManager creates a new mock instance.

func (m *MockBaseCacheManager) Add(arg0 context.Context, arg1 map[string]interface{}, arg2 interface{}) error
    Add mocks base method.

func (m *MockBaseCacheManager) Close() error
    Close mocks base method.

func (m *MockBaseCacheManager) EXPECT() *MockBaseCacheManagerMockRecorder
    EXPECT returns an object that allows the caller to indicate expected use.

func (m *MockBaseCacheManager) Flush(arg0 context.Context) error
    Flush mocks base method.

func (m *MockBaseCacheManager) Get(arg0 context.Context, arg1 map[string]interface{}) ([]byte, error)
    Get mocks base method.

func (m *MockBaseCacheManager) Has(arg0 context.Context, arg1 map[string]interface{}) (bool, error)
    Has mocks base method.

type MockBaseCacheManagerMockRecorder struct {
	// Has unexported fields.
}
    MockBaseCacheManagerMockRecorder is the mock recorder for
    MockBaseCacheManager.

func (mr *MockBaseCacheManagerMockRecorder) Add(arg0, arg1, arg2 interface{}) *gomock.Call
    Add indicates an expected call of Add.

func (mr *MockBaseCacheManagerMockRecorder) Close() *gomock.Call
    Close indicates an expected call of Close.

func (mr *MockBaseCacheManagerMockRecorder) Flush(arg0 interface{}) *gomock.Call
    Flush indicates an expected call of Flush.

func (mr *MockBaseCacheManagerMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call
    Get indicates an expected call of Get.

func (mr *MockBaseCacheManagerMockRecorder) Has(arg0, arg1 interface{}) *gomock.Call
    Has indicates an expected call of Has.

type MockBaseClient struct {
	// Has unexported fields.
}
    MockBaseClient is a mock of BaseClient interface.

func NewMockBaseClient(ctrl *gomock.Controller) *MockBaseClient
    NewMockBaseClient creates a new mock instance.

func (m *MockBaseClient) DoGetURIWithQuery(arg0 context.Context, arg1 string, arg2 map[string]string) (int64, []byte, error)
    DoGetURIWithQuery mocks base method.

func (m *MockBaseClient) EXPECT() *MockBaseClientMockRecorder
    EXPECT returns an object that allows the caller to indicate expected use.

type MockBaseClientMockRecorder struct {
	// Has unexported fields.
}
    MockBaseClientMockRecorder is the mock recorder for MockBaseClient.

func (mr *MockBaseClientMockRecorder) DoGetURIWithQuery(arg0, arg1, arg2 interface{}) *gomock.Call
    DoGetURIWithQuery indicates an expected call of DoGetURIWithQuery.

type MockBaseStorage struct {
	// Has unexported fields.
}
    MockBaseStorage is a mock of BaseStorage interface.

func NewMockBaseStorage(ctrl *gomock.Controller) *MockBaseStorage
    NewMockBaseStorage creates a new mock instance.

func (m *MockBaseStorage) AddPerson(arg0 context.Context, arg1 *datastructs.PersonData) (int64, error)
    AddPerson mocks base method.

func (m *MockBaseStorage) DeletePersonById(arg0 context.Context, arg1 int64) (*datastructs.PersonData, error)
    DeletePersonById mocks base method.

func (m *MockBaseStorage) EXPECT() *MockBaseStorageMockRecorder
    EXPECT returns an object that allows the caller to indicate expected use.

func (m *MockBaseStorage) SelectPersons(arg0 context.Context, arg1 map[string]interface{}) ([]datastructs.PersonData, error)
    SelectPersons mocks base method.

func (m *MockBaseStorage) UpdatePersonById(arg0 context.Context, arg1 int64, arg2 map[string]interface{}) (*datastructs.PersonData, error)
    UpdatePersonById mocks base method.

type MockBaseStorageManager struct {
	// Has unexported fields.
}
    MockBaseStorageManager is a mock of BaseStorageManager interface.

func NewMockBaseStorageManager(ctrl *gomock.Controller) *MockBaseStorageManager
    NewMockBaseStorageManager creates a new mock instance.

func (m *MockBaseStorageManager) AddPerson(arg0 context.Context, arg1 *datastructs.PersonData) (int64, error)
    AddPerson mocks base method.

func (m *MockBaseStorageManager) CheckConnection(arg0 context.Context) (bool, error)
    CheckConnection mocks base method.

func (m *MockBaseStorageManager) Close() error
    Close mocks base method.

func (m *MockBaseStorageManager) DeletePersonById(arg0 context.Context, arg1 int64) (int64, error)
    DeletePersonById mocks base method.

func (m *MockBaseStorageManager) EXPECT() *MockBaseStorageManagerMockRecorder
    EXPECT returns an object that allows the caller to indicate expected use.

func (m *MockBaseStorageManager) SelectPersons(arg0 context.Context, arg1 map[string]interface{}, arg2, arg3 int64) ([]datastructs.PersonData, error)
    SelectPersons mocks base method.

func (m *MockBaseStorageManager) UpdatePersonById(arg0 context.Context, arg1 int64, arg2 map[string]interface{}) (int64, error)
    UpdatePersonById mocks base method.

type MockBaseStorageManagerMockRecorder struct {
	// Has unexported fields.
}
    MockBaseStorageManagerMockRecorder is the mock recorder for
    MockBaseStorageManager.

func (mr *MockBaseStorageManagerMockRecorder) AddPerson(arg0, arg1 interface{}) *gomock.Call
    AddPerson indicates an expected call of AddPerson.

func (mr *MockBaseStorageManagerMockRecorder) CheckConnection(arg0 interface{}) *gomock.Call
    CheckConnection indicates an expected call of CheckConnection.

func (mr *MockBaseStorageManagerMockRecorder) Close() *gomock.Call
    Close indicates an expected call of Close.

func (mr *MockBaseStorageManagerMockRecorder) DeletePersonById(arg0, arg1 interface{}) *gomock.Call
    DeletePersonById indicates an expected call of DeletePersonById.

func (mr *MockBaseStorageManagerMockRecorder) SelectPersons(arg0, arg1, arg2, arg3 interface{}) *gomock.Call
    SelectPersons indicates an expected call of SelectPersons.

func (mr *MockBaseStorageManagerMockRecorder) UpdatePersonById(arg0, arg1, arg2 interface{}) *gomock.Call
    UpdatePersonById indicates an expected call of UpdatePersonById.

type MockBaseStorageMockRecorder struct {
	// Has unexported fields.
}
    MockBaseStorageMockRecorder is the mock recorder for MockBaseStorage.

func (mr *MockBaseStorageMockRecorder) AddPerson(arg0, arg1 interface{}) *gomock.Call
    AddPerson indicates an expected call of AddPerson.

func (mr *MockBaseStorageMockRecorder) DeletePersonById(arg0, arg1 interface{}) *gomock.Call
    DeletePersonById indicates an expected call of DeletePersonById.

func (mr *MockBaseStorageMockRecorder) SelectPersons(arg0, arg1 interface{}) *gomock.Call
    SelectPersons indicates an expected call of SelectPersons.

func (mr *MockBaseStorageMockRecorder) UpdatePersonById(arg0, arg1, arg2 interface{}) *gomock.Call
    UpdatePersonById indicates an expected call of UpdatePersonById.

type MockConsumer struct {
	// Has unexported fields.
}
    MockConsumer is a mock of Consumer interface.

func NewMockConsumer(ctrl *gomock.Controller) *MockConsumer
    NewMockConsumer creates a new mock instance.

func (m *MockConsumer) Close() error
    Close mocks base method.

func (m *MockConsumer) EXPECT() *MockConsumerMockRecorder
    EXPECT returns an object that allows the caller to indicate expected use.

func (m *MockConsumer) Listen(arg0 context.Context, arg1 chan<- msgbroker.Message)
    Listen mocks base method.

func (m *MockConsumer) ReadMessage(arg0 context.Context) (msgbroker.Message, error)
    ReadMessage mocks base method.

type MockConsumerMockRecorder struct {
	// Has unexported fields.
}
    MockConsumerMockRecorder is the mock recorder for MockConsumer.

func (mr *MockConsumerMockRecorder) Close() *gomock.Call
    Close indicates an expected call of Close.

func (mr *MockConsumerMockRecorder) Listen(arg0, arg1 interface{}) *gomock.Call
    Listen indicates an expected call of Listen.

func (mr *MockConsumerMockRecorder) ReadMessage(arg0 interface{}) *gomock.Call
    ReadMessage indicates an expected call of ReadMessage.

type MockProducer struct {
	// Has unexported fields.
}
    MockProducer is a mock of Producer interface.

func NewMockProducer(ctrl *gomock.Controller) *MockProducer
    NewMockProducer creates a new mock instance.

func (m *MockProducer) Close() error
    Close mocks base method.

func (m *MockProducer) EXPECT() *MockProducerMockRecorder
    EXPECT returns an object that allows the caller to indicate expected use.

func (m *MockProducer) Listen(arg0 context.Context, arg1 <-chan msgbroker.Message)
    Listen mocks base method.

func (m *MockProducer) SendMessage(arg0 context.Context, arg1, arg2 string) error
    SendMessage mocks base method.

type MockProducerMockRecorder struct {
	// Has unexported fields.
}
    MockProducerMockRecorder is the mock recorder for MockProducer.

func (mr *MockProducerMockRecorder) Close() *gomock.Call
    Close indicates an expected call of Close.

func (mr *MockProducerMockRecorder) Listen(arg0, arg1 interface{}) *gomock.Call
    Listen indicates an expected call of Listen.

func (mr *MockProducerMockRecorder) SendMessage(arg0, arg1, arg2 interface{}) *gomock.Call
    SendMessage indicates an expected call of SendMessage.


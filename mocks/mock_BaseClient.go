// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/erupshis/effective_mobile/internal/client (interfaces: BaseClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBaseClient is a mock of BaseClient interface.
type MockBaseClient struct {
	ctrl     *gomock.Controller
	recorder *MockBaseClientMockRecorder
}

// MockBaseClientMockRecorder is the mock recorder for MockBaseClient.
type MockBaseClientMockRecorder struct {
	mock *MockBaseClient
}

// NewMockBaseClient creates a new mock instance.
func NewMockBaseClient(ctrl *gomock.Controller) *MockBaseClient {
	mock := &MockBaseClient{ctrl: ctrl}
	mock.recorder = &MockBaseClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBaseClient) EXPECT() *MockBaseClientMockRecorder {
	return m.recorder
}

// DoGetURIWithQuery mocks base method.
func (m *MockBaseClient) DoGetURIWithQuery(arg0 context.Context, arg1 string, arg2 map[string]string) (int64, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DoGetURIWithQuery", arg0, arg1, arg2)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DoGetURIWithQuery indicates an expected call of DoGetURIWithQuery.
func (mr *MockBaseClientMockRecorder) DoGetURIWithQuery(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoGetURIWithQuery", reflect.TypeOf((*MockBaseClient)(nil).DoGetURIWithQuery), arg0, arg1, arg2)
}

// makeEmptyBodyRequest mocks base method.
func (m *MockBaseClient) makeEmptyBodyRequest(arg0 context.Context, arg1, arg2 string) (int64, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "makeEmptyBodyRequest", arg0, arg1, arg2)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// makeEmptyBodyRequest indicates an expected call of makeEmptyBodyRequest.
func (mr *MockBaseClientMockRecorder) makeEmptyBodyRequest(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "makeEmptyBodyRequest", reflect.TypeOf((*MockBaseClient)(nil).makeEmptyBodyRequest), arg0, arg1, arg2)
}

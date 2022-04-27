// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	io "io"
	reflect "reflect"

	static "github.com/zitadel/zitadel/internal/static"
	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// GetObject mocks base method.
func (m *MockStorage) GetObject(ctx context.Context, instanceID, resourceOwner, name string) ([]byte, func() (*static.Asset, error), error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObject", ctx, instanceID, resourceOwner, name)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(func() (*static.Asset, error))
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetObject indicates an expected call of GetObject.
func (mr *MockStorageMockRecorder) GetObject(ctx, instanceID, resourceOwner, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObject", reflect.TypeOf((*MockStorage)(nil).GetObject), ctx, instanceID, resourceOwner, name)
}

// GetObjectInfo mocks base method.
func (m *MockStorage) GetObjectInfo(ctx context.Context, instanceID, resourceOwner, name string) (*static.Asset, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObjectInfo", ctx, instanceID, resourceOwner, name)
	ret0, _ := ret[0].(*static.Asset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObjectInfo indicates an expected call of GetObjectInfo.
func (mr *MockStorageMockRecorder) GetObjectInfo(ctx, instanceID, resourceOwner, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObjectInfo", reflect.TypeOf((*MockStorage)(nil).GetObjectInfo), ctx, instanceID, resourceOwner, name)
}

// PutObject mocks base method.
func (m *MockStorage) PutObject(ctx context.Context, instanceID, location, resourceOwner, name, contentType string, objectType static.ObjectType, object io.Reader, objectSize int64) (*static.Asset, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutObject", ctx, instanceID, location, resourceOwner, name, contentType, objectType, object, objectSize)
	ret0, _ := ret[0].(*static.Asset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutObject indicates an expected call of PutObject.
func (mr *MockStorageMockRecorder) PutObject(ctx, instanceID, location, resourceOwner, name, contentType, objectType, object, objectSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutObject", reflect.TypeOf((*MockStorage)(nil).PutObject), ctx, instanceID, location, resourceOwner, name, contentType, objectType, object, objectSize)
}

// RemoveObject mocks base method.
func (m *MockStorage) RemoveObject(ctx context.Context, instanceID, resourceOwner, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveObject", ctx, instanceID, resourceOwner, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveObject indicates an expected call of RemoveObject.
func (mr *MockStorageMockRecorder) RemoveObject(ctx, instanceID, resourceOwner, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveObject", reflect.TypeOf((*MockStorage)(nil).RemoveObject), ctx, instanceID, resourceOwner, name)
}

// RemoveObjects mocks base method.
func (m *MockStorage) RemoveObjects(ctx context.Context, instanceID, resourceOwner string, objectType static.ObjectType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveObjects", ctx, instanceID, resourceOwner, objectType)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveObjects indicates an expected call of RemoveObjects.
func (mr *MockStorageMockRecorder) RemoveObjects(ctx, instanceID, resourceOwner, objectType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveObjects", reflect.TypeOf((*MockStorage)(nil).RemoveObjects), ctx, instanceID, resourceOwner, objectType)
}

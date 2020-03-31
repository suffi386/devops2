// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/caos/zitadel/internal/eventstore (interfaces: App)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	models "github.com/caos/zitadel/internal/eventstore/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockApp is a mock of App interface
type MockApp struct {
	ctrl     *gomock.Controller
	recorder *MockAppMockRecorder
}

// MockAppMockRecorder is the mock recorder for MockApp
type MockAppMockRecorder struct {
	mock *MockApp
}

// NewMockApp creates a new mock instance
func NewMockApp(ctrl *gomock.Controller) *MockApp {
	mock := &MockApp{ctrl: ctrl}
	mock.recorder = &MockAppMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockApp) EXPECT() *MockAppMockRecorder {
	return m.recorder
}

// CreateEvents mocks base method
func (m *MockApp) CreateEvents(arg0 context.Context, arg1 ...*models.Aggregate) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateEvents", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateEvents indicates an expected call of CreateEvents
func (mr *MockAppMockRecorder) CreateEvents(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEvents", reflect.TypeOf((*MockApp)(nil).CreateEvents), varargs...)
}

// FilterEvents mocks base method
func (m *MockApp) FilterEvents(arg0 context.Context, arg1 *models.SearchQuery) ([]*models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterEvents", arg0, arg1)
	ret0, _ := ret[0].([]*models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FilterEvents indicates an expected call of FilterEvents
func (mr *MockAppMockRecorder) FilterEvents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterEvents", reflect.TypeOf((*MockApp)(nil).FilterEvents), arg0, arg1)
}

// Health mocks base method
func (m *MockApp) Health(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Health", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Health indicates an expected call of Health
func (mr *MockAppMockRecorder) Health(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Health", reflect.TypeOf((*MockApp)(nil).Health), arg0)
}

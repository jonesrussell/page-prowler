// Code generated by MockGen. DO NOT EDIT.
// Source: internal/prowlredis/prowlredis.go

// Package prowlredis is a generated GoMock package.
package prowlredis

import (
	context "context"
	reflect "reflect"
	"time"

	gomock "github.com/golang/mock/gomock"
)

// MockClientInterface is a mock of ClientInterface interface.
type MockClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockClientInterfaceMockRecorder
}

// Set implements ClientInterface.
func (m *MockClientInterface) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	panic("unimplemented")
}

// MockClientInterfaceMockRecorder is the mock recorder for MockClientInterface.
type MockClientInterfaceMockRecorder struct {
	mock *MockClientInterface
}

// NewMockClientInterface creates a new mock instance.
func NewMockClientInterface(ctrl *gomock.Controller) *MockClientInterface {
	mock := &MockClientInterface{ctrl: ctrl}
	mock.recorder = &MockClientInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientInterface) EXPECT() *MockClientInterfaceMockRecorder {
	return m.recorder
}

// Del mocks base method.
func (m *MockClientInterface) Del(ctx context.Context, keys ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range keys {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Del", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Del indicates an expected call of Del.
func (mr *MockClientInterfaceMockRecorder) Del(ctx interface{}, keys ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, keys...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Del", reflect.TypeOf((*MockClientInterface)(nil).Del), varargs...)
}

// Options mocks base method.
func (m *MockClientInterface) Options() *Options {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Options")
	ret0, _ := ret[0].(*Options)
	return ret0
}

// Options indicates an expected call of Options.
func (mr *MockClientInterfaceMockRecorder) Options() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Options", reflect.TypeOf((*MockClientInterface)(nil).Options))
}

// Ping mocks base method.
func (m *MockClientInterface) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockClientInterfaceMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockClientInterface)(nil).Ping), ctx)
}

// SAdd mocks base method.
func (m *MockClientInterface) SAdd(ctx context.Context, key string, members ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, key}
	for _, a := range members {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SAdd", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SAdd indicates an expected call of SAdd.
func (mr *MockClientInterfaceMockRecorder) SAdd(ctx, key interface{}, members ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, key}, members...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SAdd", reflect.TypeOf((*MockClientInterface)(nil).SAdd), varargs...)
}

// SIsMember mocks base method.
func (m *MockClientInterface) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SIsMember", ctx, key, member)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SIsMember indicates an expected call of SIsMember.
func (mr *MockClientInterfaceMockRecorder) SIsMember(ctx, key, member interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SIsMember", reflect.TypeOf((*MockClientInterface)(nil).SIsMember), ctx, key, member)
}

// SMembers mocks base method.
func (m *MockClientInterface) SMembers(ctx context.Context, key string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SMembers", ctx, key)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SMembers indicates an expected call of SMembers.
func (mr *MockClientInterfaceMockRecorder) SMembers(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SMembers", reflect.TypeOf((*MockClientInterface)(nil).SMembers), ctx, key)
}

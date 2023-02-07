// Code generated by MockGen. DO NOT EDIT.
// Source: C:\Users\gilsp\go\src\aqua-farm-manager\pkg\redis\redis.go

// Package mock_redis is a generated GoMock package.
package mock_redis

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRedisMethod is a mock of RedisMethod interface.
type MockRedisMethod struct {
	ctrl     *gomock.Controller
	recorder *MockRedisMethodMockRecorder
}

// MockRedisMethodMockRecorder is the mock recorder for MockRedisMethod.
type MockRedisMethodMockRecorder struct {
	mock *MockRedisMethod
}

// NewMockRedisMethod creates a new mock instance.
func NewMockRedisMethod(ctrl *gomock.Controller) *MockRedisMethod {
	mock := &MockRedisMethod{ctrl: ctrl}
	mock.recorder = &MockRedisMethodMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRedisMethod) EXPECT() *MockRedisMethodMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockRedisMethod) Delete(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRedisMethodMockRecorder) Delete(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRedisMethod)(nil).Delete), key)
}

// Get mocks base method.
func (m *MockRedisMethod) Get(key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRedisMethodMockRecorder) Get(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRedisMethod)(nil).Get), key)
}

// HGETALL mocks base method.
func (m *MockRedisMethod) HGETALL(key string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HGETALL", key)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HGETALL indicates an expected call of HGETALL.
func (mr *MockRedisMethodMockRecorder) HGETALL(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HGETALL", reflect.TypeOf((*MockRedisMethod)(nil).HGETALL), key)
}

// HINCRBY mocks base method.
func (m *MockRedisMethod) HINCRBY(key, field string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HINCRBY", key, field)
	ret0, _ := ret[0].(error)
	return ret0
}

// HINCRBY indicates an expected call of HINCRBY.
func (mr *MockRedisMethodMockRecorder) HINCRBY(key, field interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HINCRBY", reflect.TypeOf((*MockRedisMethod)(nil).HINCRBY), key, field)
}

// HSET mocks base method.
func (m *MockRedisMethod) HSET(key, field, value string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HSET", key, field, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// HSET indicates an expected call of HSET.
func (mr *MockRedisMethodMockRecorder) HSET(key, field, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HSET", reflect.TypeOf((*MockRedisMethod)(nil).HSET), key, field, value)
}

// SETNX mocks base method.
func (m *MockRedisMethod) SETNX(key string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SETNX", key)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SETNX indicates an expected call of SETNX.
func (mr *MockRedisMethodMockRecorder) SETNX(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SETNX", reflect.TypeOf((*MockRedisMethod)(nil).SETNX), key)
}

// Set mocks base method.
func (m *MockRedisMethod) Set(key, field string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", key, field)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockRedisMethodMockRecorder) Set(key, field interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockRedisMethod)(nil).Set), key, field)
}

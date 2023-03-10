// Code generated by MockGen. DO NOT EDIT.
// Source: C:\Users\gilsp\go\src\aqua-farm-manager\internal\domain\farm\farm.go

// Package mock_farm is a generated GoMock package.
package mock_farm

import (
	farm "aqua-farm-manager/internal/domain/farm"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFarmDomain is a mock of FarmDomain interface.
type MockFarmDomain struct {
	ctrl     *gomock.Controller
	recorder *MockFarmDomainMockRecorder
}

// MockFarmDomainMockRecorder is the mock recorder for MockFarmDomain.
type MockFarmDomainMockRecorder struct {
	mock *MockFarmDomain
}

// NewMockFarmDomain creates a new mock instance.
func NewMockFarmDomain(ctrl *gomock.Controller) *MockFarmDomain {
	mock := &MockFarmDomain{ctrl: ctrl}
	mock.recorder = &MockFarmDomainMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFarmDomain) EXPECT() *MockFarmDomainMockRecorder {
	return m.recorder
}

// CreateFarmInfo mocks base method.
func (m *MockFarmDomain) CreateFarmInfo(r farm.CreateDomainRequest) (farm.CreateDomainResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFarmInfo", r)
	ret0, _ := ret[0].(farm.CreateDomainResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFarmInfo indicates an expected call of CreateFarmInfo.
func (mr *MockFarmDomainMockRecorder) CreateFarmInfo(r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFarmInfo", reflect.TypeOf((*MockFarmDomain)(nil).CreateFarmInfo), r)
}

// DeleteFarmInfo mocks base method.
func (m *MockFarmDomain) DeleteFarmInfo(r farm.DeleteDomainRequest) (farm.DeleteDomainResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFarmInfo", r)
	ret0, _ := ret[0].(farm.DeleteDomainResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFarmInfo indicates an expected call of DeleteFarmInfo.
func (mr *MockFarmDomainMockRecorder) DeleteFarmInfo(r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFarmInfo", reflect.TypeOf((*MockFarmDomain)(nil).DeleteFarmInfo), r)
}

// DeleteFarmsWithDependencies mocks base method.
func (m *MockFarmDomain) DeleteFarmsWithDependencies(ID uint) (farm.DeleteAllResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFarmsWithDependencies", ID)
	ret0, _ := ret[0].(farm.DeleteAllResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFarmsWithDependencies indicates an expected call of DeleteFarmsWithDependencies.
func (mr *MockFarmDomainMockRecorder) DeleteFarmsWithDependencies(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFarmsWithDependencies", reflect.TypeOf((*MockFarmDomain)(nil).DeleteFarmsWithDependencies), ID)
}

// GetFarm mocks base method.
func (m *MockFarmDomain) GetFarm(size, cursor int) ([]farm.GetFarmInfoResponse, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFarm", size, cursor)
	ret0, _ := ret[0].([]farm.GetFarmInfoResponse)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetFarm indicates an expected call of GetFarm.
func (mr *MockFarmDomainMockRecorder) GetFarm(size, cursor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFarm", reflect.TypeOf((*MockFarmDomain)(nil).GetFarm), size, cursor)
}

// GetFarmInfoByID mocks base method.
func (m *MockFarmDomain) GetFarmInfoByID(ID uint) (farm.GetFarmInfoResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFarmInfoByID", ID)
	ret0, _ := ret[0].(farm.GetFarmInfoResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFarmInfoByID indicates an expected call of GetFarmInfoByID.
func (mr *MockFarmDomainMockRecorder) GetFarmInfoByID(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFarmInfoByID", reflect.TypeOf((*MockFarmDomain)(nil).GetFarmInfoByID), ID)
}

// UpdateFarmInfo mocks base method.
func (m *MockFarmDomain) UpdateFarmInfo(r farm.UpdateDomainRequest) (farm.UpdateDomainResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFarmInfo", r)
	ret0, _ := ret[0].(farm.UpdateDomainResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateFarmInfo indicates an expected call of UpdateFarmInfo.
func (mr *MockFarmDomainMockRecorder) UpdateFarmInfo(r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFarmInfo", reflect.TypeOf((*MockFarmDomain)(nil).UpdateFarmInfo), r)
}

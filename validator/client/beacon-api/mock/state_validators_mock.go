// Code generated by MockGen. DO NOT EDIT.
// Source: validator/client/beacon-api/state_validators.go
//
// Generated by this command:
//
//	mockgen -package=mock -source=validator/client/beacon-api/state_validators.go -destination=validator/client/beacon-api/mock/state_validators_mock.go
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	structs "github.com/prysmaticlabs/prysm/v5/api/server/structs"
	primitives "github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
	gomock "go.uber.org/mock/gomock"
)

// MockStateValidatorsProvider is a mock of StateValidatorsProvider interface.
type MockStateValidatorsProvider struct {
	ctrl     *gomock.Controller
	recorder *MockStateValidatorsProviderMockRecorder
}

// MockStateValidatorsProviderMockRecorder is the mock recorder for MockStateValidatorsProvider.
type MockStateValidatorsProviderMockRecorder struct {
	mock *MockStateValidatorsProvider
}

// NewMockStateValidatorsProvider creates a new mock instance.
func NewMockStateValidatorsProvider(ctrl *gomock.Controller) *MockStateValidatorsProvider {
	mock := &MockStateValidatorsProvider{ctrl: ctrl}
	mock.recorder = &MockStateValidatorsProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStateValidatorsProvider) EXPECT() *MockStateValidatorsProviderMockRecorder {
	return m.recorder
}

// GetStateValidators mocks base method.
func (m *MockStateValidatorsProvider) GetStateValidators(arg0 context.Context, arg1 []string, arg2 []primitives.ValidatorIndex, arg3 []string) (*structs.GetValidatorsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStateValidators", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*structs.GetValidatorsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStateValidators indicates an expected call of GetStateValidators.
func (mr *MockStateValidatorsProviderMockRecorder) GetStateValidators(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStateValidators", reflect.TypeOf((*MockStateValidatorsProvider)(nil).GetStateValidators), arg0, arg1, arg2, arg3)
}

// GetStateValidatorsForHead mocks base method.
func (m *MockStateValidatorsProvider) GetStateValidatorsForHead(arg0 context.Context, arg1 []string, arg2 []primitives.ValidatorIndex, arg3 []string) (*structs.GetValidatorsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStateValidatorsForHead", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*structs.GetValidatorsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStateValidatorsForHead indicates an expected call of GetStateValidatorsForHead.
func (mr *MockStateValidatorsProviderMockRecorder) GetStateValidatorsForHead(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStateValidatorsForHead", reflect.TypeOf((*MockStateValidatorsProvider)(nil).GetStateValidatorsForHead), arg0, arg1, arg2, arg3)
}

// GetStateValidatorsForSlot mocks base method.
func (m *MockStateValidatorsProvider) GetStateValidatorsForSlot(arg0 context.Context, arg1 primitives.Slot, arg2 []string, arg3 []primitives.ValidatorIndex, arg4 []string) (*structs.GetValidatorsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStateValidatorsForSlot", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*structs.GetValidatorsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStateValidatorsForSlot indicates an expected call of GetStateValidatorsForSlot.
func (mr *MockStateValidatorsProviderMockRecorder) GetStateValidatorsForSlot(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStateValidatorsForSlot", reflect.TypeOf((*MockStateValidatorsProvider)(nil).GetStateValidatorsForSlot), arg0, arg1, arg2, arg3, arg4)
}

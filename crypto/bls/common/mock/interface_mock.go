// Code generated by MockGen. DO NOT EDIT.
// Source: crypto/bls/common/interface.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	common "github.com/prysmaticlabs/prysm/v4/crypto/bls/common"
)

// MockSecretKey is a mock of SecretKey interface.
type SecretKey struct {
	ctrl     *gomock.Controller
	recorder *SecretKeyMockRecorder
}

// MockSecretKeyMockRecorder is the mock recorder for MockSecretKey.
type SecretKeyMockRecorder struct {
	mock *SecretKey
}

// NewSecretKey creates a new mock instance.
func NewSecretKey(ctrl *gomock.Controller) *SecretKey {
	mock := &SecretKey{ctrl: ctrl}
	mock.recorder = &SecretKeyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *SecretKey) EXPECT() *SecretKeyMockRecorder {
	return m.recorder
}

// Marshal mocks base method.
func (m *SecretKey) Marshal() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Marshal")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// Marshal indicates an expected call of Marshal.
func (mr *SecretKeyMockRecorder) Marshal() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Marshal", reflect.TypeOf((*SecretKey)(nil).Marshal))
}

// PublicKey mocks base method.
func (m *SecretKey) PublicKey() common.PublicKey {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublicKey")
	ret0, _ := ret[0].(common.PublicKey)
	return ret0
}

// PublicKey indicates an expected call of PublicKey.
func (mr *SecretKeyMockRecorder) PublicKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublicKey", reflect.TypeOf((*SecretKey)(nil).PublicKey))
}

// Sign mocks base method.
func (m *SecretKey) Sign(msg []byte) common.Signature {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sign", msg)
	ret0, _ := ret[0].(common.Signature)
	return ret0
}

// Sign indicates an expected call of Sign.
func (mr *SecretKeyMockRecorder) Sign(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*SecretKey)(nil).Sign), msg)
}

// MockPublicKey is a mock of PublicKey interface.
type PublicKey struct {
	ctrl     *gomock.Controller
	recorder *PublicKeyMockRecorder
}

// MockPublicKeyMockRecorder is the mock recorder for MockPublicKey.
type PublicKeyMockRecorder struct {
	mock *PublicKey
}

// NewPublicKey creates a new mock instance.
func NewPublicKey(ctrl *gomock.Controller) *PublicKey {
	mock := &PublicKey{ctrl: ctrl}
	mock.recorder = &PublicKeyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *PublicKey) EXPECT() *PublicKeyMockRecorder {
	return m.recorder
}

// Aggregate mocks base method.
func (m *PublicKey) Aggregate(p2 common.PublicKey) common.PublicKey {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Aggregate", p2)
	ret0, _ := ret[0].(common.PublicKey)
	return ret0
}

// Aggregate indicates an expected call of Aggregate.
func (mr *PublicKeyMockRecorder) Aggregate(p2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Aggregate", reflect.TypeOf((*PublicKey)(nil).Aggregate), p2)
}

// Copy mocks base method.
func (m *PublicKey) Copy() common.PublicKey {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Copy")
	ret0, _ := ret[0].(common.PublicKey)
	return ret0
}

// Copy indicates an expected call of Copy.
func (mr *PublicKeyMockRecorder) Copy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Copy", reflect.TypeOf((*PublicKey)(nil).Copy))
}

// Equals mocks base method.
func (m *PublicKey) Equals(p2 common.PublicKey) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Equals", p2)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Equals indicates an expected call of Equals.
func (mr *PublicKeyMockRecorder) Equals(p2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Equals", reflect.TypeOf((*PublicKey)(nil).Equals), p2)
}

// IsInfinite mocks base method.
func (m *PublicKey) IsInfinite() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsInfinite")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsInfinite indicates an expected call of IsInfinite.
func (mr *PublicKeyMockRecorder) IsInfinite() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsInfinite", reflect.TypeOf((*PublicKey)(nil).IsInfinite))
}

// Marshal mocks base method.
func (m *PublicKey) Marshal() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Marshal")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// Marshal indicates an expected call of Marshal.
func (mr *PublicKeyMockRecorder) Marshal() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Marshal", reflect.TypeOf((*PublicKey)(nil).Marshal))
}

// MockSignature is a mock of Signature interface.
type Signature struct {
	ctrl     *gomock.Controller
	recorder *SignatureMockRecorder
}

// MockSignatureMockRecorder is the mock recorder for MockSignature.
type SignatureMockRecorder struct {
	mock *Signature
}

// NewSignature creates a new mock instance.
func NewSignature(ctrl *gomock.Controller) *Signature {
	mock := &Signature{ctrl: ctrl}
	mock.recorder = &SignatureMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Signature) EXPECT() *SignatureMockRecorder {
	return m.recorder
}

// AggregateVerify mocks base method.
func (m *Signature) AggregateVerify(pubKeys []common.PublicKey, msgs [][32]byte) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AggregateVerify", pubKeys, msgs)
	ret0, _ := ret[0].(bool)
	return ret0
}

// AggregateVerify indicates an expected call of AggregateVerify.
func (mr *SignatureMockRecorder) AggregateVerify(pubKeys, msgs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AggregateVerify", reflect.TypeOf((*Signature)(nil).AggregateVerify), pubKeys, msgs)
}

// Copy mocks base method.
func (m *Signature) Copy() common.Signature {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Copy")
	ret0, _ := ret[0].(common.Signature)
	return ret0
}

// Copy indicates an expected call of Copy.
func (mr *SignatureMockRecorder) Copy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Copy", reflect.TypeOf((*Signature)(nil).Copy))
}

// Eth2FastAggregateVerify mocks base method.
func (m *Signature) Eth2FastAggregateVerify(pubKeys []common.PublicKey, msg [32]byte) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Eth2FastAggregateVerify", pubKeys, msg)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Eth2FastAggregateVerify indicates an expected call of Eth2FastAggregateVerify.
func (mr *SignatureMockRecorder) Eth2FastAggregateVerify(pubKeys, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Eth2FastAggregateVerify", reflect.TypeOf((*Signature)(nil).Eth2FastAggregateVerify), pubKeys, msg)
}

// FastAggregateVerify mocks base method.
func (m *Signature) FastAggregateVerify(pubKeys []common.PublicKey, msg [32]byte) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FastAggregateVerify", pubKeys, msg)
	ret0, _ := ret[0].(bool)
	return ret0
}

// FastAggregateVerify indicates an expected call of FastAggregateVerify.
func (mr *SignatureMockRecorder) FastAggregateVerify(pubKeys, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FastAggregateVerify", reflect.TypeOf((*Signature)(nil).FastAggregateVerify), pubKeys, msg)
}

// Marshal mocks base method.
func (m *Signature) Marshal() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Marshal")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// Marshal indicates an expected call of Marshal.
func (mr *SignatureMockRecorder) Marshal() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Marshal", reflect.TypeOf((*Signature)(nil).Marshal))
}

// Verify mocks base method.
func (m *Signature) Verify(pubKey common.PublicKey, msg []byte) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", pubKey, msg)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Verify indicates an expected call of Verify.
func (mr *SignatureMockRecorder) Verify(pubKey, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*Signature)(nil).Verify), pubKey, msg)
}

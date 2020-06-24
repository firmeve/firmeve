// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/firmeve/firmeve/kernel/contract (interfaces: Render)

// Package mock is a generated GoMock package.
package mock

import (
	contract "github.com/firmeve/firmeve/kernel/contract"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRender is a mock of Render interface
type MockRender struct {
	ctrl     *gomock.Controller
	recorder *MockRenderMockRecorder
}

// MockRenderMockRecorder is the mock recorder for MockRender
type MockRenderMockRecorder struct {
	mock *MockRender
}

// NewMockRender creates a new mock instance
func NewMockRender(ctrl *gomock.Controller) *MockRender {
	mock := &MockRender{ctrl: ctrl}
	mock.recorder = &MockRenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRender) EXPECT() *MockRenderMockRecorder {
	return m.recorder
}

// Render mocks base method
func (m *MockRender) Render(arg0 contract.Protocol, arg1 int, arg2 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Render", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Render indicates an expected call of Render
func (mr *MockRenderMockRecorder) Render(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Render", reflect.TypeOf((*MockRender)(nil).Render), arg0, arg1, arg2)
}

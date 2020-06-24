// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/firmeve/firmeve/kernel/contract (interfaces: Context)

// Package mock is a generated GoMock package.
package mock

import (
	contract "github.com/firmeve/firmeve/kernel/contract"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockContext is a mock of Context interface
type MockContext struct {
	ctrl     *gomock.Controller
	recorder *MockContextMockRecorder
}

// MockContextMockRecorder is the mock recorder for MockContext
type MockContextMockRecorder struct {
	mock *MockContext
}

// NewMockContext creates a new mock instance
func NewMockContext(ctrl *gomock.Controller) *MockContext {
	mock := &MockContext{ctrl: ctrl}
	mock.recorder = &MockContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockContext) EXPECT() *MockContextMockRecorder {
	return m.recorder
}

// Abort mocks base method
func (m *MockContext) Abort() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Abort")
}

// Abort indicates an expected call of Abort
func (mr *MockContextMockRecorder) Abort() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Abort", reflect.TypeOf((*MockContext)(nil).Abort))
}

// AddEntity mocks base method
func (m *MockContext) AddEntity(arg0 string, arg1 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddEntity", arg0, arg1)
}

// AddEntity indicates an expected call of AddEntity
func (mr *MockContextMockRecorder) AddEntity(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEntity", reflect.TypeOf((*MockContext)(nil).AddEntity), arg0, arg1)
}

// Application mocks base method
func (m *MockContext) Application() contract.Application {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Application")
	ret0, _ := ret[0].(contract.Application)
	return ret0
}

// Application indicates an expected call of Application
func (mr *MockContextMockRecorder) Application() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Application", reflect.TypeOf((*MockContext)(nil).Application))
}

// Bind mocks base method
func (m *MockContext) Bind(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bind", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Bind indicates an expected call of Bind
func (mr *MockContextMockRecorder) Bind(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bind", reflect.TypeOf((*MockContext)(nil).Bind), arg0)
}

// BindValidate mocks base method
func (m *MockContext) BindValidate(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BindValidate", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// BindValidate indicates an expected call of BindValidate
func (mr *MockContextMockRecorder) BindValidate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BindValidate", reflect.TypeOf((*MockContext)(nil).BindValidate), arg0)
}

// BindWith mocks base method
func (m *MockContext) BindWith(arg0 contract.Binding, arg1 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BindWith", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// BindWith indicates an expected call of BindWith
func (mr *MockContextMockRecorder) BindWith(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BindWith", reflect.TypeOf((*MockContext)(nil).BindWith), arg0, arg1)
}

// BindWithValidate mocks base method
func (m *MockContext) BindWithValidate(arg0 contract.Binding, arg1 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BindWithValidate", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// BindWithValidate indicates an expected call of BindWithValidate
func (mr *MockContextMockRecorder) BindWithValidate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BindWithValidate", reflect.TypeOf((*MockContext)(nil).BindWithValidate), arg0, arg1)
}

// Clone mocks base method
func (m *MockContext) Clone() contract.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Clone")
	ret0, _ := ret[0].(contract.Context)
	return ret0
}

// Clone indicates an expected call of Clone
func (mr *MockContextMockRecorder) Clone() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clone", reflect.TypeOf((*MockContext)(nil).Clone))
}

// Current mocks base method
func (m *MockContext) Current() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Current")
	ret0, _ := ret[0].(int)
	return ret0
}

// Current indicates an expected call of Current
func (mr *MockContextMockRecorder) Current() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Current", reflect.TypeOf((*MockContext)(nil).Current))
}

// Deadline mocks base method
func (m *MockContext) Deadline() (time.Time, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Deadline")
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Deadline indicates an expected call of Deadline
func (mr *MockContextMockRecorder) Deadline() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deadline", reflect.TypeOf((*MockContext)(nil).Deadline))
}

// Done mocks base method
func (m *MockContext) Done() <-chan struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Done")
	ret0, _ := ret[0].(<-chan struct{})
	return ret0
}

// Done indicates an expected call of Done
func (mr *MockContextMockRecorder) Done() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Done", reflect.TypeOf((*MockContext)(nil).Done))
}

// Entity mocks base method
func (m *MockContext) Entity(arg0 string) *contract.ContextEntity {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Entity", arg0)
	ret0, _ := ret[0].(*contract.ContextEntity)
	return ret0
}

// Entity indicates an expected call of Entity
func (mr *MockContextMockRecorder) Entity(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Entity", reflect.TypeOf((*MockContext)(nil).Entity), arg0)
}

// Err mocks base method
func (m *MockContext) Err() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Err")
	ret0, _ := ret[0].(error)
	return ret0
}

// Err indicates an expected call of Err
func (mr *MockContextMockRecorder) Err() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Err", reflect.TypeOf((*MockContext)(nil).Err))
}

// Error mocks base method
func (m *MockContext) Error(arg0 int, arg1 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Error", arg0, arg1)
}

// Error indicates an expected call of Error
func (mr *MockContextMockRecorder) Error(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockContext)(nil).Error), arg0, arg1)
}

// Firmeve mocks base method
func (m *MockContext) Firmeve() contract.Application {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Firmeve")
	ret0, _ := ret[0].(contract.Application)
	return ret0
}

// Firmeve indicates an expected call of Firmeve
func (mr *MockContextMockRecorder) Firmeve() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Firmeve", reflect.TypeOf((*MockContext)(nil).Firmeve))
}

// Get mocks base method
func (m *MockContext) Get(arg0 string) interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockContextMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockContext)(nil).Get), arg0)
}

// Handlers mocks base method
func (m *MockContext) Handlers() []contract.ContextHandler {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handlers")
	ret0, _ := ret[0].([]contract.ContextHandler)
	return ret0
}

// Handlers indicates an expected call of Handlers
func (mr *MockContextMockRecorder) Handlers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handlers", reflect.TypeOf((*MockContext)(nil).Handlers))
}

// Next mocks base method
func (m *MockContext) Next() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Next")
}

// Next indicates an expected call of Next
func (mr *MockContextMockRecorder) Next() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*MockContext)(nil).Next))
}

// Protocol mocks base method
func (m *MockContext) Protocol() contract.Protocol {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Protocol")
	ret0, _ := ret[0].(contract.Protocol)
	return ret0
}

// Protocol indicates an expected call of Protocol
func (mr *MockContextMockRecorder) Protocol() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Protocol", reflect.TypeOf((*MockContext)(nil).Protocol))
}

// Render mocks base method
func (m *MockContext) Render(arg0 int, arg1 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Render", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Render indicates an expected call of Render
func (mr *MockContextMockRecorder) Render(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Render", reflect.TypeOf((*MockContext)(nil).Render), arg0, arg1)
}

// RenderWith mocks base method
func (m *MockContext) RenderWith(arg0 int, arg1 contract.Render, arg2 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RenderWith", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RenderWith indicates an expected call of RenderWith
func (mr *MockContextMockRecorder) RenderWith(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderWith", reflect.TypeOf((*MockContext)(nil).RenderWith), arg0, arg1, arg2)
}

// Resolve mocks base method
func (m *MockContext) Resolve(arg0 interface{}, arg1 ...interface{}) interface{} {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Resolve", varargs...)
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Resolve indicates an expected call of Resolve
func (mr *MockContextMockRecorder) Resolve(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resolve", reflect.TypeOf((*MockContext)(nil).Resolve), varargs...)
}

// Value mocks base method
func (m *MockContext) Value(arg0 interface{}) interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Value", arg0)
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Value indicates an expected call of Value
func (mr *MockContextMockRecorder) Value(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Value", reflect.TypeOf((*MockContext)(nil).Value), arg0)
}

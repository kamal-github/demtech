// Code generated by MockGen. DO NOT EDIT.
// Source: internal/validator/quotavalidator.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLastNHoursCountGetter is a mock of LastNHoursCountGetter interface.
type MockLastNHoursCountGetter struct {
	ctrl     *gomock.Controller
	recorder *MockLastNHoursCountGetterMockRecorder
}

// MockLastNHoursCountGetterMockRecorder is the mock recorder for MockLastNHoursCountGetter.
type MockLastNHoursCountGetterMockRecorder struct {
	mock *MockLastNHoursCountGetter
}

// NewMockLastNHoursCountGetter creates a new mock instance.
func NewMockLastNHoursCountGetter(ctrl *gomock.Controller) *MockLastNHoursCountGetter {
	mock := &MockLastNHoursCountGetter{ctrl: ctrl}
	mock.recorder = &MockLastNHoursCountGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLastNHoursCountGetter) EXPECT() *MockLastNHoursCountGetterMockRecorder {
	return m.recorder
}

// GetLastNHoursCount mocks base method.
func (m *MockLastNHoursCountGetter) GetLastNHoursCount(ctx context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastNHoursCount", ctx)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastNHoursCount indicates an expected call of GetLastNHoursCount.
func (mr *MockLastNHoursCountGetterMockRecorder) GetLastNHoursCount(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastNHoursCount", reflect.TypeOf((*MockLastNHoursCountGetter)(nil).GetLastNHoursCount), ctx)
}

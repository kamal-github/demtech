package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSentEmailTracker is a mock of SentEmailTracker interface.
type MockSentEmailTracker struct {
	ctrl     *gomock.Controller
	recorder *MockSentEmailTrackerMockRecorder
}

// MockSentEmailTrackerMockRecorder is the mock recorder for MockSentEmailTracker.
type MockSentEmailTrackerMockRecorder struct {
	mock *MockSentEmailTracker
}

// NewMockSentEmailTracker creates a new mock instance.
func NewMockSentEmailTracker(ctrl *gomock.Controller) *MockSentEmailTracker {
	mock := &MockSentEmailTracker{ctrl: ctrl}
	mock.recorder = &MockSentEmailTrackerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSentEmailTracker) EXPECT() *MockSentEmailTrackerMockRecorder {
	return m.recorder
}

// TrackSentEmail mocks base method.
func (m *MockSentEmailTracker) TrackSentEmail(ctx context.Context, msgID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TrackSentEmail", ctx, msgID)
	ret0, _ := ret[0].(error)
	return ret0
}

// TrackSentEmail indicates an expected call of TrackSentEmail.
func (mr *MockSentEmailTrackerMockRecorder) TrackSentEmail(ctx, msgID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TrackSentEmail", reflect.TypeOf((*MockSentEmailTracker)(nil).TrackSentEmail), ctx, msgID)
}

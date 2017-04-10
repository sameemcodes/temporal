package history

import "github.com/stretchr/testify/mock"
import gohistory "github.com/uber/cadence/.gen/go/history"
import "github.com/uber/cadence/.gen/go/shared"

// MockHistoryEngine is used as mock implementation for HistoryEngine
type MockHistoryEngine struct {
	mock.Mock
}

// Start is mock implementation for Start for HistoryEngine
func (_m *MockHistoryEngine) Start() {
	_m.Called()
}

// Stop is mock implementation for Stop of HistoryEngine
func (_m *MockHistoryEngine) Stop() {
	_m.Called()
}

// StartWorkflowExecution is mock implementation for StartWorkflowExecution of HistoryEngine
func (_m *MockHistoryEngine) StartWorkflowExecution(request *gohistory.StartWorkflowExecutionRequest) (*shared.StartWorkflowExecutionResponse, error) {
	ret := _m.Called(request)

	var r0 *shared.StartWorkflowExecutionResponse
	if rf, ok := ret.Get(0).(func(*gohistory.StartWorkflowExecutionRequest) *shared.StartWorkflowExecutionResponse); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*shared.StartWorkflowExecutionResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gohistory.StartWorkflowExecutionRequest) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWorkflowExecutionHistory is mock implementation for GetWorkflowExecutionHistory of HistoryEngine
func (_m *MockHistoryEngine) GetWorkflowExecutionHistory(request *gohistory.GetWorkflowExecutionHistoryRequest) (*shared.GetWorkflowExecutionHistoryResponse, error) {
	ret := _m.Called(request)

	var r0 *shared.GetWorkflowExecutionHistoryResponse
	if rf, ok := ret.Get(0).(func(*gohistory.GetWorkflowExecutionHistoryRequest) *shared.GetWorkflowExecutionHistoryResponse); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*shared.GetWorkflowExecutionHistoryResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gohistory.GetWorkflowExecutionHistoryRequest) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RecordDecisionTaskStarted is mock implementation for RecordDecisionTaskStarted of HistoryEngine
func (_m *MockHistoryEngine) RecordDecisionTaskStarted(request *gohistory.RecordDecisionTaskStartedRequest) (*gohistory.RecordDecisionTaskStartedResponse, error) {
	ret := _m.Called(request)

	var r0 *gohistory.RecordDecisionTaskStartedResponse
	if rf, ok := ret.Get(0).(func(*gohistory.RecordDecisionTaskStartedRequest) *gohistory.RecordDecisionTaskStartedResponse); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gohistory.RecordDecisionTaskStartedResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gohistory.RecordDecisionTaskStartedRequest) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RecordActivityTaskStarted is mock implementation for RecordActivityTaskStarted of HistoryEngine
func (_m *MockHistoryEngine) RecordActivityTaskStarted(request *gohistory.RecordActivityTaskStartedRequest) (*gohistory.RecordActivityTaskStartedResponse, error) {
	ret := _m.Called(request)

	var r0 *gohistory.RecordActivityTaskStartedResponse
	if rf, ok := ret.Get(0).(func(*gohistory.RecordActivityTaskStartedRequest) *gohistory.RecordActivityTaskStartedResponse); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gohistory.RecordActivityTaskStartedResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gohistory.RecordActivityTaskStartedRequest) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RespondDecisionTaskCompleted is mock implementation for RespondDecisionTaskCompleted of HistoryEngine
func (_m *MockHistoryEngine) RespondDecisionTaskCompleted(request *gohistory.RespondDecisionTaskCompletedRequest) error {
	ret := _m.Called(request)

	var r0 error
	if rf, ok := ret.Get(0).(func(*gohistory.RespondDecisionTaskCompletedRequest) error); ok {
		r0 = rf(request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RespondActivityTaskCompleted is mock implementation for RespondActivityTaskCompleted of HistoryEngine
func (_m *MockHistoryEngine) RespondActivityTaskCompleted(request *gohistory.RespondActivityTaskCompletedRequest) error {
	ret := _m.Called(request)

	var r0 error
	if rf, ok := ret.Get(0).(func(*gohistory.RespondActivityTaskCompletedRequest) error); ok {
		r0 = rf(request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RespondActivityTaskFailed is mock implementation for RespondActivityTaskFailed of HistoryEngine
func (_m *MockHistoryEngine) RespondActivityTaskFailed(request *gohistory.RespondActivityTaskFailedRequest) error {
	ret := _m.Called(request)

	var r0 error
	if rf, ok := ret.Get(0).(func(*gohistory.RespondActivityTaskFailedRequest) error); ok {
		r0 = rf(request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RespondActivityTaskCanceled is mock implementation for RespondActivityTaskCanceled of HistoryEngine
func (_m *MockHistoryEngine) RespondActivityTaskCanceled(request *gohistory.RespondActivityTaskCanceledRequest) error {
	ret := _m.Called(request)

	var r0 error
	if rf, ok := ret.Get(0).(func(*gohistory.RespondActivityTaskCanceledRequest) error); ok {
		r0 = rf(request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RecordActivityTaskHeartbeat is mock implementation for RecordActivityTaskHeartbeat of HistoryEngine
func (_m *MockHistoryEngine) RecordActivityTaskHeartbeat(request *gohistory.RecordActivityTaskHeartbeatRequest) (*shared.RecordActivityTaskHeartbeatResponse, error) {
	ret := _m.Called(request)

	var r0 *shared.RecordActivityTaskHeartbeatResponse
	if rf, ok := ret.Get(0).(func(*gohistory.RecordActivityTaskHeartbeatRequest) *shared.RecordActivityTaskHeartbeatResponse); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*shared.RecordActivityTaskHeartbeatResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gohistory.RecordActivityTaskHeartbeatRequest) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TerminateWorkflowExecution is mock implementation for TerminateWorkflowExecution of HistoryEngine
func (_m *MockHistoryEngine) TerminateWorkflowExecution(request *gohistory.TerminateWorkflowExecutionRequest) error {
	ret := _m.Called(request)

	var r0 error
	if rf, ok := ret.Get(0).(func(*gohistory.TerminateWorkflowExecutionRequest) error); ok {
		r0 = rf(request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

var _ Engine = (*MockHistoryEngine)(nil)

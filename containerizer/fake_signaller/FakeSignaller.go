// This file was generated by counterfeiter
package fake_signaller

import (
	"sync"

	"code.cloudfoundry.org/garden-linux/containerizer"
)

type FakeSignaller struct {
	SignalErrorStub        func(err error) error
	signalErrorMutex       sync.RWMutex
	signalErrorArgsForCall []struct {
		err error
	}
	signalErrorReturns struct {
		result1 error
	}
	SignalSuccessStub        func() error
	signalSuccessMutex       sync.RWMutex
	signalSuccessArgsForCall []struct{}
	signalSuccessReturns     struct {
		result1 error
	}
}

func (fake *FakeSignaller) SignalError(err error) error {
	fake.signalErrorMutex.Lock()
	fake.signalErrorArgsForCall = append(fake.signalErrorArgsForCall, struct {
		err error
	}{err})
	fake.signalErrorMutex.Unlock()
	if fake.SignalErrorStub != nil {
		return fake.SignalErrorStub(err)
	} else {
		return fake.signalErrorReturns.result1
	}
}

func (fake *FakeSignaller) SignalErrorCallCount() int {
	fake.signalErrorMutex.RLock()
	defer fake.signalErrorMutex.RUnlock()
	return len(fake.signalErrorArgsForCall)
}

func (fake *FakeSignaller) SignalErrorArgsForCall(i int) error {
	fake.signalErrorMutex.RLock()
	defer fake.signalErrorMutex.RUnlock()
	return fake.signalErrorArgsForCall[i].err
}

func (fake *FakeSignaller) SignalErrorReturns(result1 error) {
	fake.SignalErrorStub = nil
	fake.signalErrorReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeSignaller) SignalSuccess() error {
	fake.signalSuccessMutex.Lock()
	fake.signalSuccessArgsForCall = append(fake.signalSuccessArgsForCall, struct{}{})
	fake.signalSuccessMutex.Unlock()
	if fake.SignalSuccessStub != nil {
		return fake.SignalSuccessStub()
	} else {
		return fake.signalSuccessReturns.result1
	}
}

func (fake *FakeSignaller) SignalSuccessCallCount() int {
	fake.signalSuccessMutex.RLock()
	defer fake.signalSuccessMutex.RUnlock()
	return len(fake.signalSuccessArgsForCall)
}

func (fake *FakeSignaller) SignalSuccessReturns(result1 error) {
	fake.SignalSuccessStub = nil
	fake.signalSuccessReturns = struct {
		result1 error
	}{result1}
}

var _ containerizer.Signaller = new(FakeSignaller)

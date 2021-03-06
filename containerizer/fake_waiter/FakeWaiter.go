// This file was generated by counterfeiter
package fake_waiter

import (
	"sync"
	"time"

	"code.cloudfoundry.org/garden-linux/containerizer"
)

type FakeWaiter struct {
	WaitStub        func(timeout time.Duration) error
	waitMutex       sync.RWMutex
	waitArgsForCall []struct {
		timeout time.Duration
	}
	waitReturns struct {
		result1 error
	}
	IsSignalErrorStub        func(err error) bool
	isSignalErrorMutex       sync.RWMutex
	isSignalErrorArgsForCall []struct {
		err error
	}
	isSignalErrorReturns struct {
		result1 bool
	}
}

func (fake *FakeWaiter) Wait(timeout time.Duration) error {
	fake.waitMutex.Lock()
	fake.waitArgsForCall = append(fake.waitArgsForCall, struct {
		timeout time.Duration
	}{timeout})
	fake.waitMutex.Unlock()
	if fake.WaitStub != nil {
		return fake.WaitStub(timeout)
	} else {
		return fake.waitReturns.result1
	}
}

func (fake *FakeWaiter) WaitCallCount() int {
	fake.waitMutex.RLock()
	defer fake.waitMutex.RUnlock()
	return len(fake.waitArgsForCall)
}

func (fake *FakeWaiter) WaitArgsForCall(i int) time.Duration {
	fake.waitMutex.RLock()
	defer fake.waitMutex.RUnlock()
	return fake.waitArgsForCall[i].timeout
}

func (fake *FakeWaiter) WaitReturns(result1 error) {
	fake.WaitStub = nil
	fake.waitReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeWaiter) IsSignalError(err error) bool {
	fake.isSignalErrorMutex.Lock()
	fake.isSignalErrorArgsForCall = append(fake.isSignalErrorArgsForCall, struct {
		err error
	}{err})
	fake.isSignalErrorMutex.Unlock()
	if fake.IsSignalErrorStub != nil {
		return fake.IsSignalErrorStub(err)
	} else {
		return fake.isSignalErrorReturns.result1
	}
}

func (fake *FakeWaiter) IsSignalErrorCallCount() int {
	fake.isSignalErrorMutex.RLock()
	defer fake.isSignalErrorMutex.RUnlock()
	return len(fake.isSignalErrorArgsForCall)
}

func (fake *FakeWaiter) IsSignalErrorArgsForCall(i int) error {
	fake.isSignalErrorMutex.RLock()
	defer fake.isSignalErrorMutex.RUnlock()
	return fake.isSignalErrorArgsForCall[i].err
}

func (fake *FakeWaiter) IsSignalErrorReturns(result1 bool) {
	fake.IsSignalErrorStub = nil
	fake.isSignalErrorReturns = struct {
		result1 bool
	}{result1}
}

var _ containerizer.Waiter = new(FakeWaiter)

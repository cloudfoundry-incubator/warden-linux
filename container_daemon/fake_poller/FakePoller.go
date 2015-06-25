// This file was generated by counterfeiter
package fake_poller

import (
	"sync"

	"github.com/cloudfoundry-incubator/garden-linux/container_daemon"
)

type FakePoller struct {
	PollStub        func() error
	pollMutex       sync.RWMutex
	pollArgsForCall []struct{}
	pollReturns     struct {
		result1 error
	}
}

func (fake *FakePoller) Poll() error {
	fake.pollMutex.Lock()
	fake.pollArgsForCall = append(fake.pollArgsForCall, struct{}{})
	fake.pollMutex.Unlock()
	if fake.PollStub != nil {
		return fake.PollStub()
	} else {
		return fake.pollReturns.result1
	}
}

func (fake *FakePoller) PollCallCount() int {
	fake.pollMutex.RLock()
	defer fake.pollMutex.RUnlock()
	return len(fake.pollArgsForCall)
}

func (fake *FakePoller) PollReturns(result1 error) {
	fake.PollStub = nil
	fake.pollReturns = struct {
		result1 error
	}{result1}
}

var _ container_daemon.Poller = new(FakePoller)
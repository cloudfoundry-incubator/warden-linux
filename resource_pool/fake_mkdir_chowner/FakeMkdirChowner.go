// This file was generated by counterfeiter
package fake_mkdir_chowner

import (
	"os"
	"sync"

	"github.com/cloudfoundry-incubator/garden-linux/resource_pool"
)

type FakeMkdirChowner struct {
	MkdirChownStub        func(path string, uid, gid uint32, mode os.FileMode) error
	mkdirChownMutex       sync.RWMutex
	mkdirChownArgsForCall []struct {
		path string
		uid  uint32
		gid  uint32
		mode os.FileMode
	}
	mkdirChownReturns struct {
		result1 error
	}
}

func (fake *FakeMkdirChowner) MkdirChown(path string, uid uint32, gid uint32, mode os.FileMode) error {
	fake.mkdirChownMutex.Lock()
	fake.mkdirChownArgsForCall = append(fake.mkdirChownArgsForCall, struct {
		path string
		uid  uint32
		gid  uint32
		mode os.FileMode
	}{path, uid, gid, mode})
	fake.mkdirChownMutex.Unlock()
	if fake.MkdirChownStub != nil {
		return fake.MkdirChownStub(path, uid, gid, mode)
	} else {
		return fake.mkdirChownReturns.result1
	}
}

func (fake *FakeMkdirChowner) MkdirChownCallCount() int {
	fake.mkdirChownMutex.RLock()
	defer fake.mkdirChownMutex.RUnlock()
	return len(fake.mkdirChownArgsForCall)
}

func (fake *FakeMkdirChowner) MkdirChownArgsForCall(i int) (string, uint32, uint32, os.FileMode) {
	fake.mkdirChownMutex.RLock()
	defer fake.mkdirChownMutex.RUnlock()
	return fake.mkdirChownArgsForCall[i].path, fake.mkdirChownArgsForCall[i].uid, fake.mkdirChownArgsForCall[i].gid, fake.mkdirChownArgsForCall[i].mode
}

func (fake *FakeMkdirChowner) MkdirChownReturns(result1 error) {
	fake.MkdirChownStub = nil
	fake.mkdirChownReturns = struct {
		result1 error
	}{result1}
}

var _ resource_pool.MkdirChowner = new(FakeMkdirChowner)
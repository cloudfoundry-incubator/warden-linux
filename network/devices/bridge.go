package devices

import (
	"fmt"
	"net"
	"sync"

	"github.com/docker/libcontainer/netlink"
)

// netlink is not thread-safe, all calls to netlink should be guarded by this mutex
var netlinkMu *sync.Mutex = new(sync.Mutex)

type Bridge struct{}

// Create creates a bridge device and returns the interface.
// If the device already exists, returns the existing interface.
func (Bridge) Create(name string, ip net.IP, subnet *net.IPNet) (intf *net.Interface, err error) {
	netlinkMu.Lock()
	defer netlinkMu.Unlock()

	if intf, err = idempotentlyCreateBridge(name); err != nil {
		return nil, err
	}

	if err = netlink.NetworkLinkAddIp(intf, ip, subnet); err != nil && err.Error() != "file exists" {
		return nil, fmt.Errorf("devices: add IP to bridge: %v", err)
	}

	return intf, nil
}

func idempotentlyCreateBridge(name string) (intf *net.Interface, err error) {
	createErr := netlink.CreateBridge(name, true)

	intfs, listErr := net.Interfaces()
	if listErr != nil {
		return nil, fmt.Errorf("devices: list bridges: %s", listErr)
	}

	if intf, ok := findBridgeIntf(intfs, name); ok {
		return intf, nil
	}

	return nil, fmt.Errorf("devices: create bridge: ", createErr)
}

func findBridgeIntf(intfs []net.Interface, name string) (intf *net.Interface, found bool) {
	for _, intf := range intfs {
		if intf.Name == name {
			return &intf, true
		}
	}

	return nil, false
}

func (Bridge) Add(bridge, slave *net.Interface) error {
	netlinkMu.Lock()
	defer netlinkMu.Unlock()

	return netlink.AddToBridge(slave, bridge)
}

func (Bridge) Destroy(bridge string) error {
	netlinkMu.Lock()
	defer netlinkMu.Unlock()

	intfs, err := net.Interfaces()
	if err != nil {
		return err
	}

	for _, i := range intfs {
		if i.Name == bridge {
			return netlink.DeleteBridge(bridge)
		}
	}

	return nil
}

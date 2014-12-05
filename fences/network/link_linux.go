package network

import (
	"net"

	"github.com/docker/libcontainer/netlink"
)

type Link struct{}

func (Link) AddIP(intf *net.Interface, ip net.IP, subnet *net.IPNet) error {
	return netlink.NetworkLinkAddIp(intf, ip, subnet)
}

func (Link) AddDefaultGW(intf *net.Interface, ip net.IP) error {
	return netlink.AddDefaultGw(ip.String(), intf.Name)
}

func (Link) SetUp(intf *net.Interface) error {
	return netlink.NetworkLinkUp(intf)
}

func (Link) SetMTU(intf *net.Interface, mtu int) error {
	return netlink.NetworkSetMTU(intf, mtu)
}

func (Link) SetNs(intf *net.Interface, ns int) error {
	return netlink.NetworkSetNsPid(intf, ns)
}

func (Link) InterfaceByName(name string) (*net.Interface, bool, error) {
	intfs, err := net.Interfaces()
	if err != nil {
		return nil, false, err
	}

	for _, intf := range intfs {
		if intf.Name == name {
			return &intf, true, nil
		}
	}

	return nil, false, nil
}

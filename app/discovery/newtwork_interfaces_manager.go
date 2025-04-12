package discovery

import (
	"net"
	"sync"
)

type IPAddress net.IPNet

func NewIpAddressFromIPNet(ipNet net.IPNet) IPAddress {
	return IPAddress(ipNet)
}

func (ipNet IPAddress) getBroadcastAddress() net.IP {
	ip := ipNet.IP.To4()
	mask := ipNet.Mask

	broadcast := make(net.IP, 4)
	for i := 0; i < 4; i++ {
		broadcast[i] = ip[i] | ^mask[i]
	}
	return broadcast
}

type NetworkInterfaceManager struct {
	AvailableIpInterface []IPAddress
	InterfacesMutex      sync.Mutex
}

func NewNetworkInterfaceManager() *NetworkInterfaceManager {
	return &NetworkInterfaceManager{
		AvailableIpInterface: make([]IPAddress, 0),
	}
}

func (networkInterfaceManager *NetworkInterfaceManager) fetchInterfaces() {
	networkInterfaceManager.InterfacesMutex.Lock()
	defer networkInterfaceManager.InterfacesMutex.Unlock()

	networkInterfaceManager.AvailableIpInterface = make([]IPAddress, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok || ipnet.IP.To4() == nil {
				continue // ignore IPv6 ou adresses non-IPNet
			}
			ipAddress := NewIpAddressFromIPNet(*ipnet)
			networkInterfaceManager.AvailableIpInterface = append(networkInterfaceManager.AvailableIpInterface, ipAddress)
		}
	}
}

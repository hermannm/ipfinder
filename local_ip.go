package ipfinder

import (
	"errors"
	"fmt"
	"net"

	"hermannm.dev/wrap"
)

type LocalIP struct {
	Address          net.IP
	NetworkInterface net.Interface
}

// Goes through network interfaces on your computer to find your local IP addresses.
// Returns a list of the found addresses along with their associated network interface.
func FindLocalIPs() ([]LocalIP, error) {
	networkInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	var localIPs []LocalIP
	var errs []error

	for _, networkInterface := range networkInterfaces {
		addresses, err := networkInterface.Addrs()
		if err != nil {
			errs = append(errs, fmt.Errorf(
				"failed to get addresses for network interface '%s': %w",
				networkInterface.Name,
				err,
			))
			continue
		}

		for _, address := range addresses {
			var ip net.IP
			switch address := address.(type) {
			case *net.IPNet:
				ip = address.IP
			case *net.IPAddr:
				ip = address.IP
			}

			// Discards invalid or non-local IPs
			if ip == nil || !ip.IsPrivate() {
				continue
			}

			localIPs = append(localIPs, LocalIP{Address: ip, NetworkInterface: networkInterface})
		}
	}

	if len(localIPs) == 0 {
		if len(errs) == 0 {
			return nil, errors.New("no valid local IPs found")
		} else {
			return nil, wrap.Errors("no valid local IPs found", errs...)
		}
	}

	return localIPs, nil
}

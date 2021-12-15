package ipfinder

import (
	"errors"
	"net"
	"strings"
)

// Goes through network interfaces on your computer, and finds local IPs.
// Returns a map of network interface names to a list of IP strings connected to that interface.
// Returns error if it failed to find valid IPs.
func FindLocalIPs() (map[string][]string, error) {
	interfaces, err := net.Interfaces()

	if err != nil {
		return nil, err
	}

	localIPs := make(map[string][]string)

	for _, interf := range interfaces {
		addresses, err := interf.Addrs()

		if err != nil {
			continue
		}

		for _, address := range addresses {
			var ip net.IP

			// Asserts address as valid IP type.
			switch v := address.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// Discards invalid or non-local IPs.
			if ip == nil || !ip.IsPrivate() {
				continue
			}

			// Removes the significant bit number, as this function only wishes to return the address in itself.
			ipString := strings.Split(ip.String(), "/")[0]

			localIPs[interf.Name] = append(localIPs[interf.Name], ipString)
		}
	}

	if len(localIPs) == 0 {
		return nil, errors.New("no local IPs found")
	}

	return localIPs, nil
}

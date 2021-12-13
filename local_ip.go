package ipfinder

import (
	"errors"
	"net"
	"strings"
)

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

			switch v := address.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || !ip.IsPrivate() {
				continue
			}

			ipString := strings.Split(ip.String(), "/")[0]

			localIPs[interf.Name] = append(localIPs[interf.Name], ipString)
		}
	}

	if len(localIPs) == 0 {
		return nil, errors.New("no local IPs found")
	}

	return localIPs, nil
}

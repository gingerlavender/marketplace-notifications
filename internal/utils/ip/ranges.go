package ip

import "net"

func IsAllowed(ip net.IP, cidrs []string) bool {
	for _, cidr := range cidrs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}

		if ipNet.Contains(ip) {
			return true
		}
	}

	return false
}

package fgip

import (
	"fmt"
	"net"
	"strings"
)

func AnonymizeIP(ip string) string {
	parsed := net.ParseIP(ip)
	if parsed.To4() != nil {
		// IPv4
		parts := strings.Split(ip, ".")
		return fmt.Sprintf("%s.%s.%s.0", parts[0], parts[1], parts[2])
	}
	// IPv6
	parts := strings.Split(ip, ":")
	return strings.Join(parts[:4], ":") + "::"
}

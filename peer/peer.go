package peer

import (
	"strings"
)

func Peer(peer string) string {
	if strings.Contains(peer, "192.168.0") || strings.Contains(peer, "192.168.1") {
		return "ISP-A"
	} else if strings.Contains(peer, "10.0.0") {
		return "ISP-B"
	} else if strings.Contains(peer, "10.0.1") {
		return "ISP-C"
	} else {
		return "unknown peer"
	}
}

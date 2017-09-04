package peer

import (
	"strings"
)

func Peer(peer string) string {
	if strings.Contains(peer, "192.168.0") {
		return "ISP-hoge"
	} else if strings.Contains(peer, "10.0.0.0") {
		return "IX-hoge"
	} else {
		return "unknown peer"
	}
}

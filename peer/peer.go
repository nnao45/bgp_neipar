package peer

import (
	"strings"
)

func Peer(peer string) string {
	if strings.Contains(peer, "157.7.36") || strings.Contains(peer, "157.7.37") {
		return "GMO_JP_OCN"
	} else if strings.Contains(peer, "133.130.12") || strings.Contains(peer, "210.172.191") {
		return "GMO_JP"
	} else if strings.Contains(peer, "133.130.4") {
		return "SUNA_RentalDNS"
	} else if strings.Contains(peer, "157.7.40") {
		return "GMO_OSK"
	} else if strings.Contains(peer, "157.7.38") {
		return "GMO_Digirock"
	} else if strings.Contains(peer, "110.44.176") {
		return "mixi_private"
	} else if strings.Contains(peer, "202.94.181") || strings.Contains(peer, "202.94.193") {
		return "SuzuyoRT"
	} else if strings.Contains(peer, "210.171.224") {
		return "JPIX-IPv4"
	} else if strings.Contains(peer, "103.246.232") {
		return "JPIX-OSK"
	} else if strings.Contains(peer, "61.120.146") {
		return "GIN"
	} else if strings.Contains(peer, "218.100.6") {
		return "BBIX"
	} else if strings.Contains(peer, "203.190.230") {
		return "Equinix"
	} else if strings.Contains(peer, "202.215.242") || strings.Contains(peer, "203.140.47") {
		return "ARTERIA"
	} else if strings.Contains(peer, "183.91.61") {
		return "China-Telecom"
	} else if strings.Contains(peer, "59.128.58") {
		return "KDDI"
	} else if strings.Contains(peer, "103.234.168") {
		return "GMO_SG"
	} else if strings.Contains(peer, "61.8.56") {
		return "PACNET"
	} else if strings.Contains(peer, "103.1.208") {
		return "VIETTEL"
	} else if strings.Contains(peer, "163.44.204") {
		return "GMO_VN"
	} else if strings.Contains(peer, "118.69.56") {
		return "FPT"
	} else if strings.Contains(peer, "150.95.28") {
		return "GMO_TH"
	} else if strings.Contains(peer, "202.183.234") {
		return "CSLOXINFO"
	} else if strings.Contains(peer, "4.79.238") {
		return "LEVEL3"
	} else if strings.Contains(peer, "211.125.92") {
		return "GMO_US"
	} else if strings.Contains(peer, "64.125.220") {
		return "ZAYO"
	} else if strings.Contains(peer, "64.124.178") {
		return "ABOVENET"
	} else {
		return "unknown peer"
	}
}

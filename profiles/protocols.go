package profiles

import "strings"

type Protocol int

const (
	INVALID Protocol = iota
	SSH              // SSH with Key
	SSHPWD           // SSH with Password
	RDP
	TELNET
	REST
	K8S
)

var stop = map[string]Protocol{
	"invalid": INVALID,
	"ssh":     SSH,
	"sshpwd":  SSHPWD,
	"rdp":     RDP,
	"telnet":  TELNET,
}

// StringToProtocal parses a string into a Protocol. Returns 0 if proto is an invalid Protocol.
func StringToProtocol(proto string) Protocol {
	if proto == "" {
		return 0
	}

	proto = strings.ToLower(proto)
	p, ok := stop[proto]
	if !ok {
		return 0
	}

	return p
}

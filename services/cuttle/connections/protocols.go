package connections

import "strings"

type Protocol int

const (
	INVALID Protocol = 0
	SSH     Protocol = 1
	RDP     Protocol = 2
	TELNET  Protocol = 3
	REST    Protocol = 4
	K8S     Protocol = 5
	MOCK    Protocol = 6
)

var (
	stop = map[string]Protocol{
		"invalid": INVALID,
		"ssh":     SSH,
		"rdp":     RDP,
		"telnet":  TELNET,
		"rest":    REST,
		"k8s":     K8S,
		"mock":    MOCK,
	}
	ptos map[Protocol]string
)

func init() {
	// Reverse the StringToProtocol map so we only have one list to maintain.
	ptos = make(map[Protocol]string, len(stop))
	for s, p := range stop {
		ptos[p] = s
	}
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

// String converts the Protocol into a string. SSH => "ssh", RDP => "rdp", etc.
func (p Protocol) String() string { return ptos[p] }

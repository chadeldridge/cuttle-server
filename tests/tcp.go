package tests

import (
	"net"
	"strconv"
	"time"

	"github.com/chadeldridge/cuttle/connections"
)

const TCPDefaultTimeout = time.Second * 3

func getArg(args []map[string]any, key string) any {
	if len(args) == 0 {
		return nil
	}

	for _, a := range args {
		if v, ok := a[key]; ok {
			return v
		}
	}

	return nil
}

func getTimeoutArg(args []map[string]any) time.Duration {
	t := getArg(args, "timeout")
	if t == nil || t.(time.Duration) == 0 {
		return TCPDefaultTimeout
	}

	return t.(time.Duration)
}

func PortOpen(server connections.Server, port int, args ...map[string]any) error {
	conn, err := net.DialTimeout(
		"tcp",
		net.JoinHostPort(server.GetHostAddr(), strconv.Itoa(port)),
		getTimeoutArg(args),
	)
	if err != nil {
		return err
	}

	conn.Close()
	return nil
}

package tests

import (
	"bytes"
	"fmt"
	"time"

	"github.com/chadeldridge/cuttle-server/services/cuttle/connections"
	probing "github.com/prometheus-community/pro-bing"
)

const (
	PingDefaultCount          = 1
	PingDefaultTimeout        = time.Second * 10
	PingDefaultSuccessPercent = 1
)

// PingTest is a struct that holds the parameters for a Ping test.
type PingTest struct {
	successPercent float32
	count          int
	timeout        time.Duration
}

// NewPingTest creates a new Ping test with the given parameters.
// name: The name of the test.
// mustSucceed: If false, the Tile will continue with the test stack if this test fails.
// successPerc: Must receive packets greater than or equal to successPerc to pass. Float 0 - 1
//
//	Defaults to 1 (100%). If 0, the test will only pass if no packets are received.
//
// These TestArg will be evaluated:
// "timeout": (int, int64, time.Duration) int/int64 will be converted into time.Second * int.
// "count": int. Number of packets to send. Default is 1.
func NewPingTest(name string, mustSucceed bool, successPerc float32, args ...TestArg) Test {
	if successPerc < 0 {
		successPerc = 0
	}

	if successPerc > 1 {
		successPerc = 1
	}

	return Test{
		Name:        name,
		MustSucceed: mustSucceed,
		Tester: &PingTest{
			successPercent: successPerc,
			count:          getPingCount(args),
			timeout:        getPingTimeout(args),
		},
	}
}

func getPingCount(args []TestArg) int {
	v := FindArg(args, "count")
	if v == nil || v.(int) == 0 {
		return PingDefaultCount
	}

	return v.(int)
}

func getPingTimeout(args []TestArg) time.Duration { return GetTimeout(args, PingDefaultTimeout) }

func (p PingTest) runPinger(pinger *probing.Pinger, bufs connections.Buffers, quiet bool) error {
	buf := &bytes.Buffer{}

	if !quiet {
		pinger.OnSend = func(pkt *probing.Packet) {
			fmt.Fprintf(buf, "%d bytes from %s: icmp_seq=%d\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq)
		}
		pinger.OnRecv = func(pkt *probing.Packet) {
			fmt.Fprintf(buf, "%d bytes from %s: icmp_seq=%d time=%v\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		}

		pinger.OnDuplicateRecv = func(pkt *probing.Packet) {
			fmt.Fprintf(buf, "%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.TTL)
		}

		pinger.OnFinish = func(stats *probing.Statistics) {
			fmt.Fprintf(buf, `--- %s ping statistics ---
%d packets transmitted, %d packets received, %v%% packet loss
round-trip min/avg/max/stddev = %v/%v/%v/%v`,
				stats.Addr, stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss,
				stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		}
	}

	err := pinger.Run()
	if err != nil {
		return err
	}

	if !quiet {
		bufs.Log(time.Now(), fmt.Sprintf("PING %s (%s):\n%s", pinger.Addr(), pinger.IPAddr(), buf.String()))
	}

	return nil
}

// Ping use a UDP ping using the pro-bing library. Returns nil if successful.
// These TestArg will be evaluated:
// "quiet": bool. If true, the output will not be printed Buffers.Logs.
func (p PingTest) Run(server connections.Server, args ...TestArg) error {
	pinger, err := probing.NewPinger(server.GetHostAddr())
	if err != nil {
		return err
	}

	pinger.Count = p.count
	pinger.Timeout = p.timeout
	err = p.runPinger(pinger, server.Buffers, BeQuiet(args))
	if err != nil {
		return err
	}

	rec := float32(pinger.Statistics().PacketsRecv / p.count)
	if p.successPercent == 0 {
		if rec > 0 {
			return ErrTestFailed
		}
	}

	if rec < p.successPercent {
		return ErrTestFailed
	}

	return nil
}

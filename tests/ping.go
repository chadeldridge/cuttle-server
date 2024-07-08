package tests

import (
	"bytes"
	"fmt"
	"time"

	"github.com/chadeldridge/cuttle/connections"
	probing "github.com/prometheus-community/pro-bing"
)

const PingDefaultCount = 1

var ErrTestFailed = fmt.Errorf("failed")

func runPinger(pinger *probing.Pinger, server connections.Server) error {
	buf := &bytes.Buffer{}

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

	err := pinger.Run()
	if err != nil {
		return err
	}
	server.Log(time.Now(), fmt.Sprintf("PING %s (%s):\n%s", pinger.Addr(), pinger.IPAddr(), buf.String()))

	return nil
}

// Ping use a UDP ping using the pro-bing library. Returns true if at least one packet was received.
func Ping(server connections.Server, count int) error {
	if count == 0 {
		count = PingDefaultCount
	}

	pinger, err := probing.NewPinger(server.GetHostAddr())
	if err != nil {
		return err
	}

	pinger.Count = count
	err = runPinger(pinger, server)
	if err != nil {
		return err
	}

	if pinger.Statistics().PacketsRecv == 0 {
		return ErrTestFailed
	}

	return nil

	/*
		res := "failed"
		stats := pinger.Statistics()
		if stats.PacketsRecv > 0 {
			res = "ok"
		}

		server.PrintResults(
			time.Now(),
			fmt.Sprintf("%s (Tx: %d, Rx: %d, Lost: %d)",
				res, stats.PacketsSent, stats.PacketsRecv, stats.PacketsSent-stats.PacketsRecv,
			),
			nil,
		)

		return nil
	*/
}

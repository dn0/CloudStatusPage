package vpc

import (
	"context"
	"fmt"
	"net"

	probing "github.com/prometheus-community/pro-bing"

	"cspage/pkg/mon/agent"
)

const (
	minPongPackets = 3
	maxPacketLoss  = 10.0 // %
)

type PingAction struct {
	pinger *probing.Pinger
	ipAddr *net.IPAddr
	host   string // For display purposes only. Example: AWS VM has a name test-vm-123, but an addr i-6a7b8c.

	Id   uint32 `json:"id"`
	Size int    `json:"size"`
}

func (pa *PingAction) String() string {
	return fmt.Sprintf("host=%s addr=%s size=%d", pa.host, pa.ipAddr.String(), pa.Size)
}

func (pa *PingAction) pingRun(ctx context.Context, cfg *agent.PingConfig, host, addr string) (*probing.Statistics, error) {
	if err := pa.pingNew(cfg, host, addr, false); err != nil {
		return nil, err
	}

	err := pa.pinger.RunWithContext(ctx)
	stats := pa.pinger.Statistics()

	if err != nil {
		err = fmt.Errorf("ping(%s).Run: %w", pa, err)
	} else if stats.PacketLoss > maxPacketLoss {
		//nolint:err113 // No need for a static error.
		err = fmt.Errorf("ping(%s).Run: %.1f%% packet loss", pa, stats.PacketLoss)
	}

	return stats, err
}

func (pa *PingAction) pingPong(ctx context.Context, cfg *agent.PingConfig, host, addr string) (*probing.Statistics, error) {
	if err := pa.pingNew(cfg, host, addr, true); err != nil {
		return nil, err
	}

	err := pa.pinger.RunWithContext(ctx)
	stats := pa.pinger.Statistics()

	if err != nil {
		err = fmt.Errorf("ping(%s).Pong: %w", pa, err)
	} else if stats.PacketsRecv == 0 {
		//nolint:err113 // No need for a static error.
		err = fmt.Errorf("ping(%s).Pong: no reply (sent %d packets)", pa, stats.PacketsSent)
	}

	return stats, err
}

func (pa *PingAction) pingNew(cfg *agent.PingConfig, host, addr string, pong bool) error {
	pinger := probing.New(addr)
	pinger.Size = pa.Size
	pinger.ResolveTimeout = cfg.ResolveTimeout
	pinger.SetNetwork("ip4")

	if pong {
		pinger.Count = -1
		pinger.Timeout = cfg.PongTimeout
		pinger.Interval = cfg.PongInterval
		pinger.RecordRtts = false
		pinger.OnRecv = func(_ *probing.Packet) {
			if pinger.PacketsRecv >= minPongPackets {
				pinger.Stop()
			}
		}
		pa.ipAddr = nil
	} else {
		pinger.Count = int(cfg.Count)
		pinger.Timeout = cfg.Timeout
		pinger.Interval = cfg.Interval
	}

	pa.host = host
	pa.pinger = pinger

	if pa.ipAddr == nil {
		return pa.resolve()
	}

	// using IP address from previous resolve
	pinger.SetIPAddr(pa.ipAddr)
	return nil
}

func (pa *PingAction) resolve() error {
	err := pa.pinger.Resolve()
	if err != nil {
		return fmt.Errorf("ping(%s).Resolve: %w", pa, err)
	}
	pa.ipAddr = pa.pinger.IPAddr()
	return nil
}

func (pa *PingAction) pingReset() {
	pa.pinger = nil
	pa.ipAddr = nil
}

package sysinfo

import (
	"context"
	"log/slog"
	gonet "net"
	"slices"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"

	"cspage/pkg/pb"
)

type OSInfo struct {
	host.InfoStat
	Interfaces net.InterfaceStatList `json:"interfaces"`
}

func (os *OSInfo) PrimaryIPAddress() string {
	for _, iface := range os.Interfaces {
		if strings.HasPrefix(iface.Name, "lo") || slices.Contains(iface.Flags, "loopback") {
			continue
		}
		for _, addr := range iface.Addrs {
			ip, _, err := gonet.ParseCIDR(addr.Addr)
			if err != nil {
				continue
			}
			if ip.To4() != nil && ip.IsPrivate() {
				slog.Debug("Found primary IP address", "addr", ip)
				return ip.String()
			}
		}
	}

	slog.Error("Primary IP address not found", "interfaces", os.Interfaces)
	return ""
}

func GetOSInfo(ctx context.Context) *OSInfo {
	hostInfo, err := host.InfoWithContext(ctx)
	if err != nil {
		slog.Error("Could not fetch host info", "err", err)
		return nil
	}

	netInfo, err := net.InterfacesWithContext(ctx)
	if err != nil {
		slog.Error("Could not fetch network interfaces", "err", err)
		return nil
	}

	return &OSInfo{
		InfoStat:   *hostInfo,
		Interfaces: netInfo,
	}
}

func getOSStats(ctx context.Context) *pb.OSStat {
	ret := &pb.OSStat{}

	if cpuStats, err := cpu.TimesWithContext(ctx, false); err == nil {
		ret.Cpu = &pb.OSStat_CPU{
			User:    float32(cpuStats[0].User),
			System:  float32(cpuStats[0].System),
			Idle:    float32(cpuStats[0].Idle),
			Nice:    float32(cpuStats[0].Nice),
			Iowait:  float32(cpuStats[0].Iowait),
			Irq:     float32(cpuStats[0].Irq),
			Softirq: float32(cpuStats[0].Softirq),
			Steal:   float32(cpuStats[0].Steal),
		}
	} else {
		slog.Warn("Could not fetch OS CPU times", "err", err)
	}

	if memStats, err := mem.VirtualMemoryWithContext(ctx); err == nil {
		ret.Mem = &pb.OSStat_Memory{
			Total:        memStats.Total,
			Available:    memStats.Available,
			Used:         memStats.Used,
			Free:         memStats.Free,
			Active:       memStats.Active,
			Inactive:     memStats.Inactive,
			Wired:        memStats.Wired,
			Laundry:      memStats.Laundry,
			Buffers:      memStats.Buffers,
			Cached:       memStats.Cached,
			WriteBack:    memStats.WriteBack,
			Dirty:        memStats.Dirty,
			WriteBackTmp: memStats.WriteBackTmp,
			Shared:       memStats.Shared,
			Slab:         memStats.Slab,
		}
	} else {
		slog.Warn("Could not fetch OS virtual memory stats", "err", err)
	}

	return ret
}

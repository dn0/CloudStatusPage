package sysinfo

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/shirou/gopsutil/v3/process"

	"cspage/pkg/pb"
)

var (
	//nolint:gochecknoglobals // proc is a singleton that represents the currently running process
	proc *process.Process
	//nolint:gochecknoglobals // must be global to support the ^proc^ singleton
	once sync.Once
)

func getProc(ctx context.Context) *process.Process {
	var err error

	if proc == nil {
		once.Do(func() {
			pid := os.Getpid()
			//nolint:gosec // This integer should never overflow.
			if proc, err = process.NewProcessWithContext(ctx, int32(pid)); err != nil {
				slog.Error("Could not get process info", "pid", pid, "err", err)
				os.Exit(1)
			}
		})
	}

	return proc
}

func getProcStats(ctx context.Context) *pb.ProcStat {
	var err error
	proc := getProc(ctx)
	ret := &pb.ProcStat{}

	if ret.Threads, err = proc.NumThreadsWithContext(ctx); err != nil {
		slog.Warn("Could not get number of threads", "pid", proc.Pid, "err", err)
	}

	if ret.Fds, err = proc.NumFDsWithContext(ctx); err != nil {
		slog.Warn("Could not get number of file descriptors", "pid", proc.Pid, "err", err)
	}

	if cpuPercent, err := proc.CPUPercentWithContext(ctx); err == nil {
		ret.CpuPercent = float32(cpuPercent)
	} else {
		slog.Warn("Could not get process CPU usage", "pid", proc.Pid, "err", err)
	}

	if mem, err := proc.MemoryInfoWithContext(ctx); err == nil {
		ret.Mem = &pb.ProcStat_Memory{
			Rss:    mem.RSS,
			Vms:    mem.VMS,
			Hwm:    mem.HWM,
			Data:   mem.Data,
			Stack:  mem.Stack,
			Locked: mem.Locked,
			Swap:   mem.Swap,
		}
	} else {
		slog.Warn("Could not get process memory stats", "pid", proc.Pid, "err", err)
	}

	if ioStat, err := proc.IOCountersWithContext(ctx); err == nil {
		ret.Io = &pb.ProcStat_IO{
			ReadCount:  ioStat.ReadCount,
			WriteCount: ioStat.WriteCount,
			ReadBytes:  ioStat.ReadBytes,
			WriteBytes: ioStat.WriteBytes,
		}
	} else {
		slog.Warn("Could not get process IO stats", "pid", proc.Pid, "err", err)
	}

	return ret
}

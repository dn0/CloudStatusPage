package sysinfo

import (
	"context"

	"cspage/pkg/pb"
)

func GetSysStat(ctx context.Context) *pb.SysStat {
	return &pb.SysStat{
		Os:   getOSStats(ctx),
		Proc: getProcStats(ctx),
	}
}

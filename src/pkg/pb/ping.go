package pb

import (
	"context"
	"fmt"

	"cspage/pkg/db"
	"cspage/pkg/msg"
)

const (
	JobNamePing   = "ping"
	sqlInsertPing = `INSERT INTO {schema}.mon_ping (
  job_id, time,
  os_mem_total, os_mem_available, os_mem_used, os_mem_free, os_mem_active, os_mem_inactive, os_mem_wired, os_mem_laundry,
  os_mem_buffers, os_mem_cached, os_mem_write_back, os_mem_dirty, os_mem_write_back_tmp, os_mem_shared, os_mem_slab,
  os_cpu_user, os_cpu_system, os_cpu_idle, os_cpu_nice, os_cpu_iowait, os_cpu_irq, os_cpu_softirq, os_cpu_steal,
  proc_threads, proc_fds, proc_cpu_percent,
  proc_mem_rss, proc_mem_vms, proc_mem_hwm, proc_mem_data, proc_mem_stack, proc_mem_locked, proc_mem_swap,
  proc_io_read_count, proc_io_write_count, proc_io_read_bytes, proc_io_write_bytes
 ) VALUES (
  $1, $2,
  $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
  $18, $19, $20, $21, $22, $23, $24, $25,
  $26, $27, $28,
  $29, $30, $31, $32, $33, $34, $35,
  $36, $37, $38, $39)`
)

func (p *Ping) Repr() string {
	return fmt.Sprintf("<Ping:%s:%s>", p.GetJob().GetName(), p.ID())
}

func (p *Ping) ID() string {
	return p.GetJob().ID()
}

//nolint:protogetter,wrapcheck // Running on a freshly marshaled Object and DB error is wrapped by caller.
func (p *Ping) Save(ctx context.Context, dbc db.Client, attrs *msg.Attrs) error {
	sqlJob, sqlJobArgs := p.Job.sqlInsert(attrs)
	batch := &db.Batch{}
	batch.Queue(sqlJob, sqlJobArgs...)
	batch.Queue(
		db.WithSchema(sqlInsertPing, attrs.Cloud),
		p.Job.Id,
		p.Job.Time.AsTime(),

		p.Sysstat.Os.Mem.GetTotal(),
		p.Sysstat.Os.Mem.GetAvailable(),
		p.Sysstat.Os.Mem.GetUsed(),
		p.Sysstat.Os.Mem.GetFree(),
		p.Sysstat.Os.Mem.GetActive(),
		p.Sysstat.Os.Mem.GetInactive(),
		p.Sysstat.Os.Mem.GetWired(),
		p.Sysstat.Os.Mem.GetLaundry(),
		p.Sysstat.Os.Mem.GetBuffers(),
		p.Sysstat.Os.Mem.GetCached(),
		p.Sysstat.Os.Mem.GetWriteBack(),
		p.Sysstat.Os.Mem.GetDirty(),
		p.Sysstat.Os.Mem.GetWriteBackTmp(),
		p.Sysstat.Os.Mem.GetShared(),
		p.Sysstat.Os.Mem.GetSlab(),

		p.Sysstat.Os.Cpu.GetUser(),
		p.Sysstat.Os.Cpu.GetSystem(),
		p.Sysstat.Os.Cpu.GetIdle(),
		p.Sysstat.Os.Cpu.GetNice(),
		p.Sysstat.Os.Cpu.GetIowait(),
		p.Sysstat.Os.Cpu.GetIrq(),
		p.Sysstat.Os.Cpu.GetSoftirq(),
		p.Sysstat.Os.Cpu.GetSteal(),

		p.Sysstat.Proc.GetThreads(),
		p.Sysstat.Proc.GetFds(),
		p.Sysstat.Proc.GetCpuPercent(),

		p.Sysstat.Proc.Mem.GetRss(),
		p.Sysstat.Proc.Mem.GetVms(),
		p.Sysstat.Proc.Mem.GetHwm(),
		p.Sysstat.Proc.Mem.GetData(),
		p.Sysstat.Proc.Mem.GetStack(),
		p.Sysstat.Proc.Mem.GetLocked(),
		p.Sysstat.Proc.Mem.GetSwap(),

		p.Sysstat.Proc.Io.GetReadCount(),
		p.Sysstat.Proc.Io.GetWriteCount(),
		p.Sysstat.Proc.Io.GetReadBytes(),
		p.Sysstat.Proc.Io.GetWriteBytes(),
	)
	return dbc.SendBatch(ctx, batch).Close()
}

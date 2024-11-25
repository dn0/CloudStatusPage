SET search_path = '${CLOUD}';

-- ${CLOUD}.mon_ping
CREATE TABLE mon_ping
(
  job_id                UUID        NOT NULL, -- FK to mon_job.id
  time                  TIMESTAMPTZ NOT NULL, -- HT
  os_mem_total          BIGINT      NOT NULL,
  os_mem_available      BIGINT      NOT NULL,
  os_mem_used           BIGINT      NOT NULL,
  os_mem_free           BIGINT      NOT NULL,
  os_mem_active         BIGINT      NOT NULL,
  os_mem_inactive       BIGINT      NOT NULL,
  os_mem_wired          BIGINT      NOT NULL,
  os_mem_laundry        BIGINT      NOT NULL,
  os_mem_buffers        BIGINT      NOT NULL,
  os_mem_cached         BIGINT      NOT NULL,
  os_mem_write_back     BIGINT      NOT NULL,
  os_mem_dirty          BIGINT      NOT NULL,
  os_mem_write_back_tmp BIGINT      NOT NULL,
  os_mem_shared         BIGINT      NOT NULL,
  os_mem_slab           BIGINT      NOT NULL,
  os_cpu_user           REAL        NOT NULL,
  os_cpu_system         REAL        NOT NULL,
  os_cpu_idle           REAL        NOT NULL,
  os_cpu_nice           REAL        NOT NULL,
  os_cpu_iowait         REAL        NOT NULL,
  os_cpu_irq            REAL        NOT NULL,
  os_cpu_softirq        REAL        NOT NULL,
  os_cpu_steal          REAL        NOT NULL,
  proc_threads          INTEGER     NOT NULL,
  proc_fds              INTEGER     NOT NULL,
  proc_cpu_percent      REAL        NOT NULL,
  proc_mem_rss          BIGINT      NOT NULL,
  proc_mem_vms          BIGINT      NOT NULL,
  proc_mem_hwm          BIGINT      NOT NULL,
  proc_mem_data         BIGINT      NOT NULL,
  proc_mem_stack        BIGINT      NOT NULL,
  proc_mem_locked       BIGINT      NOT NULL,
  proc_mem_swap         BIGINT      NOT NULL,
  proc_io_read_count    BIGINT      NOT NULL,
  proc_io_write_count   BIGINT      NOT NULL,
  proc_io_read_bytes    BIGINT      NOT NULL,
  proc_io_write_bytes   BIGINT      NOT NULL
);
SELECT public.create_hypertable('mon_ping', public.by_range('time', INTERVAL '1 day'));
CREATE INDEX mon_ping_job_id_fkey ON mon_ping (job_id);

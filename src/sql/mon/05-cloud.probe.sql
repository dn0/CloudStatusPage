SET search_path = '${CLOUD}';

-- ${CLOUD}.mon_probe_${PROBE_NAME}
CREATE TABLE "mon_probe_${PROBE_NAME}"
(
  job_id   UUID        NOT NULL, -- FK to mon_job.id
  time     TIMESTAMPTZ NOT NULL, -- HT
  action   SMALLINT    NOT NULL, -- TODO: index => used by sqlSelectProbeResults
  status   SMALLINT    NOT NULL, -- TODO: index => used by sqlSelectProbeResults
  latency  BIGINT      NOT NULL,
  error    TEXT        NOT NULL
);
SELECT public.create_hypertable('mon_probe_${PROBE_NAME}', public.by_range('time', INTERVAL '7 days'));
CREATE INDEX "mon_probe_${PROBE_NAME}_job_id_fkey" ON "mon_probe_${PROBE_NAME}" (job_id);

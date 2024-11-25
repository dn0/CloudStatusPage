SET search_path = '${CLOUD}';

-- ${CLOUD}.mon_job
CREATE TABLE mon_job
(
  agent_id UUID        NOT NULL, -- FK to mon_agent.id
  id       UUID        NOT NULL, -- PK
  time     TIMESTAMPTZ NOT NULL, -- HT
  drift    BIGINT      NOT NULL,
  took     BIGINT      NOT NULL,
  name     TEXT        NOT NULL,
  errors   SMALLINT    NOT NULL
);
SELECT public.create_hypertable('mon_job', public.by_range('time', INTERVAL '1 day'));
CREATE UNIQUE INDEX mon_job_pkey ON mon_job (id DESC, time DESC); -- because of HT we need to include 'time'
CREATE INDEX mon_job_agent_id_fkey ON mon_job (agent_id);
CREATE INDEX mon_job_name_idx ON mon_job (name);

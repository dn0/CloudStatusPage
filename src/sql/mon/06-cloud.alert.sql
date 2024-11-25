SET search_path = '${CLOUD}';

-- ${CLOUD}.mon_alert
CREATE TABLE mon_alert
(
  id           UUID        NOT NULL PRIMARY KEY,
  job_id       UUID        NOT NULL, -- "FK" to mon_job.id set during alert opening
  incident_id  UUID,                 -- "FK" to mon_incident.id
  created      TIMESTAMPTZ NOT NULL,
  updated      TIMESTAMPTZ NOT NULL,
  time_begin   TIMESTAMPTZ NOT NULL, -- "FK" to mon_probe_${PROBE_NAME}.time or mon_ping.time
  time_end     TIMESTAMPTZ,          -- "FK" to mon_probe_${PROBE_NAME}.time or mon_ping.time
  time_check   TIMESTAMPTZ NOT NULL,
  type         SMALLINT    NOT NULL, -- TODO: index? => used by sqlSelectAlerts in analyst
  status       SMALLINT    NOT NULL,
  cloud_region TEXT        NOT NULL,
  probe_name   TEXT        NOT NULL, -- TODO: index? => used by sqlSelectAlertIssue in web
  probe_action SMALLINT    NOT NULL, -- 0 = no action
  data         JSONB       NOT NULL
);
CREATE INDEX mon_alert_job_id_fkey ON mon_alert (job_id);
CREATE INDEX mon_alert_incident_id_fkey ON mon_alert (incident_id);
CREATE INDEX mon_alert_created_idx ON mon_alert (created);
CREATE INDEX mon_alert_status_idx ON mon_alert (status);
CREATE INDEX mon_alert_cloud_region_idx ON mon_alert (cloud_region);

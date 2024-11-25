SET search_path = '${CLOUD}';

-- ${CLOUD}.mon_agent
CREATE TABLE mon_agent
(
  id           UUID     NOT NULL PRIMARY KEY,
  status       SMALLINT NOT NULL,
  started      TIMESTAMPTZ,
  stopped      TIMESTAMPTZ,
  version      TEXT     NOT NULL,
  hostname     TEXT     NOT NULL,
  ip_address   TEXT     NOT NULL,
  cloud_region TEXT     NOT NULL,
  cloud_zone   TEXT     NOT NULL,
  sysinfo      JSONB    NOT NULL
);
CREATE INDEX mon_agent_cloud_region_idx ON mon_agent (cloud_region);

SET search_path = '${CLOUD}';

-- ${CLOUD}.mon_incident
CREATE TABLE mon_incident
(
  id            UUID         NOT NULL PRIMARY KEY,
  created       TIMESTAMPTZ  NOT NULL,
  updated       TIMESTAMPTZ  NOT NULL,
  time_begin    TIMESTAMPTZ  NOT NULL,
  time_end      TIMESTAMPTZ,
  severity      SMALLINT     NOT NULL,
  status        SMALLINT     NOT NULL,
  cloud_regions TEXT[]       NOT NULL,
  data          JSONB        NOT NULL
);
CREATE INDEX mon_incident_created_idx ON mon_incident (created);
CREATE INDEX mon_incident_status_idx ON mon_incident (status);
CREATE INDEX mon_incident_cloud_regions_idx ON mon_incident USING GIN(cloud_regions);

SET search_path = '${CLOUD}';

-- ${CLOUD}.mon_config_region
CREATE TABLE mon_config_region
(
  id       INTEGER NOT NULL UNIQUE, -- display order
  name     TEXT    NOT NULL PRIMARY KEY,
  location TEXT    NOT NULL,
  enabled  BOOLEAN NOT NULL,
  lat      REAL,
  lon      REAL
);

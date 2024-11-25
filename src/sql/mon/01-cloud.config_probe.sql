SET search_path = '${CLOUD}';

-- ${CLOUD}.mon_config_probe
CREATE TABLE mon_config_probe
(
  id          INTEGER NOT NULL UNIQUE, -- display order
  name        TEXT    NOT NULL PRIMARY KEY,
  description TEXT    NOT NULL,
  enabled     BOOLEAN NOT NULL,
  config      JSONB   NOT NULL
);

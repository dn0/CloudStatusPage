SET search_path = '${CLOUD}';

-- ${CLOUD}.mon_checkpoint
CREATE TABLE mon_checkpoint
(
  name TEXT        NOT NULL PRIMARY KEY,
  time TIMESTAMPTZ NOT NULL
);

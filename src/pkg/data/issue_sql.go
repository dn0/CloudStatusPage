package data

import (
	"time"

	"cspage/pkg/config"
)

const (
	sqlTableAlert       = "mon_alert"
	sqlTableIncident    = "mon_incident"
	sqlCountAlertIssues = `SELECT
'{schema}' AS "cloud_id",
COUNT(DISTINCT mon_alert.id) AS "num_alerts",
0 AS "num_incidents"
FROM {schema}.mon_alert AS mon_alert
LEFT JOIN {schema}.mon_config_region AS mon_region ON mon_region.name = mon_alert.cloud_region`
	sqlCountIncidentIssues = `SELECT
'{schema}' AS "cloud_id",
0 AS "num_alerts",
COUNT(DISTINCT mon_incident.id) AS "num_incidents"
FROM {schema}.mon_incident AS mon_incident
LEFT JOIN {schema}.mon_config_region AS mon_region ON mon_region.name = mon_incident.cloud_regions[1]`
	sqlSelectAlertIssue = `SELECT DISTINCT ON (mon_alert.id)
  mon_alert.id           AS "id",
  1                      AS "type",
  '{schema}'             AS "cloud_id",
  ARRAY[mon_alert.cloud_region] AS "cloud_regions",
  mon_alert.created      AS "created",
  mon_alert.updated      AS "updated",
  mon_alert.time_begin   AS "time_begin",
  mon_alert.time_end     AS "time_end",
  0                      AS "severity",
  mon_alert.status       AS "status",
  '{}'::JSONB            AS "data",
  mon_alert.type         AS "alert_type",
  mon_alert.probe_name   AS "alert_probe_name",
  mon_alert.probe_action AS "alert_probe_action",
  mon_alert.data         AS "alert_data",
  COALESCE(mon_alert.incident_id::TEXT, '') AS "alert_incident_id"
FROM {schema}.mon_alert AS mon_alert
  LEFT JOIN {schema}.mon_config_region AS mon_region ON mon_region.name = mon_alert.cloud_region`
	sqlSelectIncidentIssue = `SELECT DISTINCT ON (mon_incident.id)
  mon_incident.id            AS "id",
  2                          AS "type",
  '{schema}'                 AS "cloud_id",
  mon_incident.cloud_regions AS "cloud_regions",
  mon_incident.created       AS "created",
  mon_incident.updated       AS "updated",
  mon_incident.time_begin    AS "time_begin",
  mon_incident.time_end      AS "time_end",
  mon_incident.severity      AS "severity",
  mon_incident.status        AS "status",
  mon_incident.data          AS "data",
  -1 AS "alert_type",
  '' AS "alert_probe_name",
  0  AS "alert_probe_action",
  '{}'::JSONB AS "alert_data",
  '' AS "alert_incident_id"
FROM {schema}.mon_incident AS mon_incident
  LEFT JOIN {schema}.mon_config_region AS mon_region ON mon_region.name = mon_incident.cloud_regions[1]`
	sqlSelectIssuesBaseCondition = `mon_region.enabled = TRUE`
	sqlSelectIncidentJoinAlerts  = `
  LEFT JOIN {schema}.mon_alert AS mon_alert ON mon_alert.incident_id = mon_incident.id`
	sqlSelectIssuesDefaultTimeSpan = 30 * 24 * time.Hour * config.DefaultTimeSpanMultiplier
)

package data

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"cspage/pkg/db"
)

const (
	sqlSelectExtendedAlerts = `SELECT
  mon_alert.id AS "id",
  mon_alert.job_id AS "job_id",
  COALESCE(mon_alert.incident_id::TEXT, '') AS "incident_id",
  mon_alert.created AS "created",
  mon_alert.updated AS "updated",
  mon_alert.time_begin AS "time_begin",
  mon_alert.time_end AS "time_end",
  mon_alert.time_check AS "time_check",
  mon_alert.type AS "type",
  mon_alert.status AS "status",
  mon_alert.cloud_region AS "cloud_region",
  mon_alert.probe_name AS "probe_name",
  mon_alert.probe_action AS "probe_action",
  mon_alert.data as "data",
  mon_incident.updated AS "incident_updated",
  mon_region.lat AS "cloud_region_lat",
  mon_region.lon AS "cloud_region_lon"
FROM {schema}.mon_alert AS mon_alert
LEFT JOIN {schema}.mon_incident AS mon_incident ON mon_incident.id = mon_alert.incident_id
LEFT JOIN {schema}.mon_config_region AS mon_region ON mon_region.name = mon_alert.cloud_region`

	weightTimeBegin        = 0.03
	weightRegion           = 0.1
	weightType             = 0.1
	weightProbeService     = 0.3
	weightProbeActionTitle = 0.05

	ClusterMaxDistance = 1.0
	clusterInfDistance = 999.0
)

type UmbrellaIncident struct {
	Id       string     `json:"id"`
	Updated  *time.Time `json:"updated"`
	Outdated bool       `json:"outdated"`
}

type ExtendedAlert struct {
	Alert
	IncidentUpdated *time.Time `json:"incident_updated"`
	CloudRegionLat  float64    `json:"cloud_region_lat"`
	CloudRegionLon  float64    `json:"cloud_region_lon"`

	distance map[*ExtendedAlert]float64
}

type AlertCluster struct {
	Alerts   []*ExtendedAlert  `json:"alerts"`
	Incident *UmbrellaIncident `json:"incident"`
}

func (a *ExtendedAlert) GetDistance(alert *ExtendedAlert) (float64, bool) {
	val, ok := a.distance[alert]
	return val, ok
}

func GetExtendedAlerts(
	ctx context.Context,
	dbc db.Conn,
	cloud string,
	filters string,
	args ...any,
) ([]*ExtendedAlert, error) {
	query := db.WithSchema(sqlSelectExtendedAlerts, cloud)
	if filters != "" {
		query += "\nWHERE " + filters
	}
	rows, _ := dbc.Query(ctx, query, args...)
	//nolint:wrapcheck // Error is properly logged by the caller.
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[ExtendedAlert])
}

//nolint:gocognit,cyclop // ClusterAlerts groups alerts into clusters based on their similarity.
func ClusterAlerts(cloud string, alerts []*ExtendedAlert, maxDistance float64) []AlertCluster {
	var clusters []AlertCluster

	for _, alert := range alerts {
		alert.distance = make(map[*ExtendedAlert]float64)
		added := false
		for i, cluster := range clusters {
			for _, clusterAlert := range cluster.Alerts {
				distance := calculateDistance(cloud, alert, clusterAlert)
				alert.distance[clusterAlert] = distance
				clusterAlert.distance[alert] = distance
				// The 2nd condition is to keep manually attached incidents in one cluster
				if distance <= maxDistance || (alert.IncidentId != "" && alert.IncidentId == clusterAlert.IncidentId) {
					clusters[i].Alerts = append(clusters[i].Alerts, alert)
					added = true
					inc := clusters[i].Incident
					if inc.Id == "" && alert.IncidentId != "" {
						inc.Id = alert.IncidentId
						inc.Updated = alert.IncidentUpdated
					}
					if inc.Updated != nil && alert.Updated.After(*inc.Updated) {
						inc.Outdated = true
					}
					break
				}
			}
			if added {
				break
			}
		}
		if !added {
			clusters = append(clusters, AlertCluster{
				Alerts: []*ExtendedAlert{alert},
				Incident: &UmbrellaIncident{
					Id:       alert.IncidentId,
					Updated:  alert.IncidentUpdated,
					Outdated: alert.IncidentUpdated != nil && alert.Updated.After(*alert.IncidentUpdated),
				},
			})
		}
	}

	return clusters
}

//nolint:mnd // calculateDistance computes a weighted distance between two alerts.
func calculateDistance(cloud string, alert1, alert2 *ExtendedAlert) float64 {
	if isInternalAlert(cloud, alert1) || isInternalAlert(cloud, alert2) {
		return clusterInfDistance
	}

	regionDist := math.Sqrt(
		math.Pow(alert1.CloudRegionLat-alert2.CloudRegionLat, 2) +
			math.Pow(alert1.CloudRegionLon-alert2.CloudRegionLon, 2))
	timeDist := math.Abs(alert1.TimeBegin.Sub(alert2.TimeBegin).Minutes())
	typeDist := math.Abs(float64(alert1.Type) - float64(alert2.Type))
	serviceDist := boolToFloat64(alert1.Data.ServiceName != alert2.Data.ServiceName)
	actionTitleDist := boolToFloat64(alert1.Data.ProbeActionTitle != alert2.Data.ProbeActionTitle)

	return weightTimeBegin*timeDist +
		weightRegion*regionDist +
		weightType*typeDist +
		weightProbeService*serviceDist +
		weightProbeActionTitle*actionTitleDist
}

func isInternalAlert(cloud string, alert *ExtendedAlert) bool {
	switch cloud {
	case cloudAWS.Id:
		if strings.Contains(alert.Data.ProbeError, "InsufficientInstanceCapacity") {
			return true
		}
	case cloudAzure.Id:
		if strings.Contains(alert.Data.ProbeError, "ConflictingUserInput") {
			return true
		}
	case cloudGCP.Id:
		if strings.Contains(alert.Data.ProbeError, "ZONE_RESOURCE_POOL_EXHAUSTED") {
			return true
		}
	}
	return false
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

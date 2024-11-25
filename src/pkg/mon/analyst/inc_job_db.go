package analyst

import (
	"context"

	"cspage/pkg/data"
	"cspage/pkg/pb"
)

const (
	// Analyzing non-ping alerts that are either open or recently updated or attached to an open incident.
	sqlSelectRecentAlertsFilters = `(
 mon_alert.updated > $1 OR
 mon_alert.status = $2 OR
 mon_incident.status = $3
) AND mon_alert.type != $4`
)

func (j *incidentJob) getRecentAlerts(ctx context.Context) ([]*data.ExtendedAlert, error) {
	//nolint:wrapcheck // Error is properly logged by the caller.
	return data.GetExtendedAlerts(
		ctx,
		j.dbc.Read,
		j.cloud,
		sqlSelectRecentAlertsFilters,
		j.checkpoint.Add(-j.cfg.IncidentCheckWindow),
		pb.AlertStatus_ALERT_OPEN,
		pb.IncidentStatus_INCIDENT_OPEN,
		pb.AlertType_PING_MISSING,
	)
}

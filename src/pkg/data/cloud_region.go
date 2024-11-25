package data

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"cspage/pkg/db"
	"cspage/pkg/pb"
)

const (
	sqlSelectCloudRegion = `SELECT DISTINCT ON (mon_region.id)
  mon_region.name     AS "name",
  mon_region.location AS "location"
FROM {schema}.mon_config_region AS mon_region
  LEFT JOIN {schema}.mon_agent as mon_agent ON mon_agent.cloud_region = mon_region.name
WHERE mon_region.enabled = TRUE AND mon_agent.status = $1`
	sqlSelectCloudRegionGeo = `SELECT
  '{schema}'          AS "cloud_id",
  mon_region.name     AS "name",
  mon_region.location AS "location",
  mon_region.enabled  AS "enabled",
  mon_region.lat      AS "lat",
  mon_region.lon      AS "lon"
FROM {schema}.mon_config_region AS mon_region`
)

type CloudRegion struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type CloudRegionGeo struct {
	CloudRegion
	CloudId string   `json:"cloud_id"`
	Enabled bool     `json:"enabled"`
	Lat     *float32 `json:"lat"`
	Lon     *float32 `json:"lon"`
}

func (r *CloudRegion) URLPrefix() string {
	if r == nil {
		return ""
	}
	return "/region/" + r.Name
}

func (i *CloudRegionGeo) Cloud() *Cloud {
	return CloudMap[i.CloudId]
}

func GetCloudRegion(ctx context.Context, dbc db.Client, cloud, region string) (*CloudRegion, error) {
	query := db.WithSchema(sqlSelectCloudRegion, cloud) + " AND mon_region.name = $2\nORDER BY mon_region.id"
	rows, _ := dbc.Query(ctx, query, pb.AgentAction_AGENT_START, region)
	obj, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[CloudRegion])
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, &db.ObjectNotFoundError{Object: fmt.Sprintf("region=%s cloud=%s", region, cloud)}
	}
	return obj, nil
}

func GetCloudRegions(ctx context.Context, dbc db.Client, cloud string) ([]*CloudRegion, error) {
	query := db.WithSchema(sqlSelectCloudRegion, cloud) + "\nORDER BY mon_region.id"
	rows, _ := dbc.Query(ctx, query, pb.AgentAction_AGENT_START)
	//nolint:wrapcheck // Error is properly logged by the caller.
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[CloudRegion])
}

func GetCloudRegionsGeo(ctx context.Context, dbc db.Client, clouds []string) ([]*CloudRegionGeo, error) {
	queries := make([]string, len(clouds))
	for i, cloud := range clouds {
		queries[i] = db.WithSchema(sqlSelectCloudRegionGeo, cloud)
	}
	rows, _ := dbc.Query(ctx, strings.Join(queries, "\nUNION ALL\n"))
	//nolint:wrapcheck // Error is properly logged by the caller.
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[CloudRegionGeo])
}

package data

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/jackc/pgx/v5"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"cspage/pkg/config"
	"cspage/pkg/db"
	"cspage/pkg/mon/agent"
	"cspage/pkg/pb"
)

const (
	ProbeMaxDisplayActionId   uint32 = 999
	probeActionGroupLeaderMod uint32 = 10
	sqlSelectProbeDefinition         = `SELECT
  mon_config_probe.name        AS "name",
  mon_config_probe.description AS "description",
  mon_config_probe.config      AS "config"
FROM {schema}.mon_config_probe AS mon_config_probe
WHERE mon_config_probe.enabled = TRUE`

	fakePingProbeName        = pb.JobNamePing
	fakePingProbeDescription = "Monitoring Agent"

	ProbeIntervalStandard probeIntervalType = "standard"
	ProbeIntervalLong     probeIntervalType = "long"
	probeIntervalPing     probeIntervalType = "ping"

	errorMaxLength = 512
)

var (
	//nolint:gochecknoglobals // This is a constant.
	fakePingConfig = &probeConfig{
		ServiceName:  "",
		ServiceGroup: "Internal Service",
		IntervalType: probeIntervalPing,
		DisplayUnits: "µs",
	}
	//nolint:gochecknoglobals // agentConfig is a singleton that contains the default mon-agent configuration
	agentConfig *agent.Config
	//nolint:gochecknoglobals // must be global to support the ^agentConfig^ singleton
	onceAgentConfig sync.Once
)

type probeIntervalType string

type ProbeAction struct {
	Id          uint32 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	ChartStack  string `json:"chart_stack,omitempty"`
}

type probeConfig struct {
	actions      map[uint32]*ProbeAction
	actionsOnce  sync.Once
	Actions      []ProbeAction     `json:"actions"`
	ServiceName  string            `json:"service_name"`
	ServiceGroup string            `json:"service_group"`
	IntervalType probeIntervalType `json:"interval_type"`
	DisplayUnits string            `json:"display_units,omitempty"`
}

type ProbeDefinition struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Config      *probeConfig `json:"config"`
}

func (p *ProbeDefinition) URLPrefix() string {
	if p == nil {
		return ""
	}
	return "/probe/" + p.Name
}

func (p *ProbeDefinition) DetailsURL(cloud *Cloud, region *CloudRegion) string {
	return cloud.URLPrefix() + region.URLPrefix() + p.URLPrefix()
}

func (p *ProbeDefinition) ChartsURL(cloud *Cloud, region *CloudRegion, qs url.Values) string {
	u := cloud.URLPrefix() + region.URLPrefix() + p.URLPrefix() + "/charts"
	if len(qs) > 0 {
		u += "?" + qs.Encode()
	}
	return u
}

func (p *ProbeDefinition) IssuesURL(cloud *Cloud, region *CloudRegion) string {
	return cloud.URLPrefix() + region.URLPrefix() + p.URLPrefix() + IssuesURLSuffix
}

func (p *ProbeDefinition) IsPingDefinition() bool {
	return p.Name == fakePingProbeName
}

// LatencyRounding returns the rounding duration, and number of decimal places to display the value in milliseconds.
// Expected to be used with something like:
//
//	strconv.FormatFloat(float64(<latency>.Round(<rounding>)) / 1e6, 'f', <decimal_places>, 64)
func (p *ProbeDefinition) LatencyRounding() (time.Duration, int) {
	//nolint:mnd // OK, these numbers are magic, but they make graphs look nicer.
	switch p.Config.DisplayUnits {
	case "ns":
		return 10 * time.Nanosecond, 5
	case "µs":
		return 10 * time.Microsecond, 2
	default: // ms
		return time.Millisecond, 0
	}
}

// GroupId uses integer division to create Action group IDs, e.g. 31 -> 30, 20 -> 20, 15 -> 10.
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (pa *ProbeAction) GroupId() uint32 {
	return pa.Id / probeActionGroupLeaderMod * probeActionGroupLeaderMod
}

func (pa *ProbeAction) ShortName() string {
	return pa.Name[strings.LastIndex(pa.Name, ".")+1:]
}

func (pa *ProbeAction) Title() string {
	return cases.Title(language.English, cases.NoLower).String(pa.ShortName())
}

//nolint:mnd // Just a bit of magic.
func (pa *ProbeAction) groupTitle() string {
	shortName := pa.ShortName()
	if len(shortName) > 3 {
		split := strings.IndexFunc(shortName[1:], func(r rune) bool {
			return unicode.IsUpper(r) || unicode.IsDigit(r)
		})
		if split > 2 {
			shortName = shortName[:split+1]
		}
	}
	return cases.Title(language.English, cases.NoLower).String(shortName)
}

func (pa *ProbeAction) FullName(p *ProbeDefinition) string {
	return ProbeActionFullName(p.Description, pa.Name, pa.Title(), false)
}

func (pa *ProbeAction) FullGroupName(p *ProbeDefinition) string {
	return ProbeActionFullName(p.Description, pa.Name, pa.groupTitle(), false)
}

func (pc *probeConfig) Interval() time.Duration {
	switch pc.IntervalType {
	case ProbeIntervalStandard:
		return getAgentConfig().ProbeIntervalDefault
	case ProbeIntervalLong:
		return getAgentConfig().ProbeLongIntervalDefault
	case probeIntervalPing:
		return getAgentConfig().PingInterval
	}
	return 0
}

func (pc *probeConfig) Timeout() time.Duration {
	switch pc.IntervalType {
	case ProbeIntervalStandard:
		return getAgentConfig().ProbeTimeout
	case ProbeIntervalLong:
		return getAgentConfig().ProbeLongTimeout
	case probeIntervalPing:
		return getAgentConfig().WorkerTaskTimeout
	}
	return 0
}

// ActionGroupIDs returns only group Action IDs.
func (pc *probeConfig) ActionGroupIDs() []uint32 {
	var actions []uint32
	for _, pa := range pc.Actions {
		if pa.Id%probeActionGroupLeaderMod == 0 {
			actions = append(actions, pa.GroupId())
		}
	}
	return actions
}

// ActionMap returns ProbeAction.Id mapped to ProbeAction object.
func (pc *probeConfig) ActionMap() map[uint32]*ProbeAction {
	if pc.actions == nil {
		pc.actionsOnce.Do(func() {
			pc.actions = make(map[uint32]*ProbeAction, len(pc.Actions))
			for _, pa := range pc.Actions {
				pc.actions[pa.Id] = &pa
			}
		})
	}
	return pc.actions
}

func (pc *probeConfig) ActionGet(aid uint32) (*ProbeAction, bool) {
	action, ok := pc.ActionMap()[aid]
	return action, ok
}

// GroupActionIDs create groups of actions, e.g. [[10,11], [20], [30,31]].
func GroupActionIDs(actions []uint32) [][]uint32 {
	slices.Sort(actions)
	var groups [][]uint32
	i := -1
	for _, aid := range actions {
		if aid%probeActionGroupLeaderMod == 0 || i < 0 {
			groups = append(groups, []uint32{})
			i++
		}
		groups[i] = append(groups[i], aid)
	}
	return groups
}

func GetProbeDefinition(ctx context.Context, dbc db.Client, cloud, name string) (*ProbeDefinition, error) {
	if name == fakePingProbeName {
		return NewPingProbeDefinition(), nil
	}

	query := db.WithSchema(sqlSelectProbeDefinition, cloud) + " AND mon_config_probe.name = $1"
	rows, _ := dbc.Query(ctx, query, name)
	obj, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[ProbeDefinition])
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, &db.ObjectNotFoundError{Object: fmt.Sprintf("probe=%s cloud=%s", name, cloud)}
	}
	return obj, nil
}

func GetProbeDefinitions(ctx context.Context, dbc db.Client, cloud string) ([]*ProbeDefinition, error) {
	query := db.WithSchema(sqlSelectProbeDefinition, cloud)
	query += "\nORDER BY mon_config_probe.id"
	rows, _ := dbc.Query(ctx, query)
	//nolint:wrapcheck // Error is properly logged by the caller.
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[ProbeDefinition])
}

func NewPingProbeDefinition() *ProbeDefinition {
	return &ProbeDefinition{
		Name:        fakePingProbeName,
		Description: fakePingProbeDescription,
		Config:      fakePingConfig,
	}
}

func ParseProbeError(msg string) string {
	// TODO: replace IDs and UUIDs.
	replacer := strings.NewReplacer(
		"----", "", // Azure uses a lot of `---` in error messages.
	)
	msg = replacer.Replace(msg)
	if len(msg) > errorMaxLength {
		msg = msg[:errorMaxLength] + "..."
	}
	return msg
}

func ProbeActionFullName(probeDescription, actionName, actionTitle string, verbose bool) string {
	switch {
	case strings.Contains(actionName, ".inter."):
		return actionTitle + " between " + probeDescription + "s"
	case strings.Contains(actionName, ".intra."):
		return actionTitle + " in " + probeDescription
	case strings.Contains(actionName, ".latency"):
		return actionTitle + " of " + probeDescription
	case verbose:
		return actionTitle + " operation of " + probeDescription
	default:
		return actionTitle + " " + probeDescription
	}
}

func getAgentConfig() *agent.Config {
	if agentConfig == nil {
		onceAgentConfig.Do(func() {
			agentConfig = &agent.Config{}
			config.InitConfigFlags(agentConfig, "fake-")
		})
	}
	return agentConfig
}

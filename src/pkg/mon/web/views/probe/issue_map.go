package probe

import (
	"strconv"

	"cspage/pkg/data"
	"cspage/pkg/pb"
)

type simpleIssue struct {
	category data.IssueType
	severity pb.IncidentSeverity
}

type issueMap struct {
	counter data.IssueCounter
	m       map[string]*simpleIssue
}

func (i *issueMap) get(cloudRegion, probeName string) *simpleIssue {
	return i.m[issueMapKey(cloudRegion, probeName)]
}

func (i *issueMap) get2(cloudRegion, probeName string, probeAction uint32) *simpleIssue {
	item := i.m[issueMapKey(cloudRegion, probeName)]
	if item != nil && item.category == data.IssueTypeIncident {
		return item // Incidents "affect" all probe actions and take precedence over alerts
	}
	return i.m[issueMapActionKey(cloudRegion, probeName, probeAction)]
}

func (i *issueMap) add(key string, issue *data.Issue) {
	item, exists := i.m[key]
	if !exists {
		item = &simpleIssue{}
		i.m[key] = item
	}
	severity := issue.Severity()
	if item.severity < severity {
		item.severity = severity
		item.category = issue.Type
	}
}

func newIssueMap(issues []*data.Issue) *issueMap {
	imap := &issueMap{m: make(map[string]*simpleIssue)}
	for _, issue := range issues {
		imap.counter.Inc(issue.Type)
		for _, service := range issue.AffectedServices() {
			if service.ProbeName == "" {
				continue
			}
			imap.add(issueMapKey(service.CloudRegion, service.ProbeName), issue)
			if issue.AlertProbeAction != 0 {
				imap.add(issueMapActionKey(service.CloudRegion, service.ProbeName, issue.AlertProbeAction), issue)
			}
		}
	}
	return imap
}

func issueMapKey(cloudRegion, probeName string) string {
	return cloudRegion + "|" + probeName
}

func issueMapActionKey(cloudRegion, probeName string, probeAction uint32) string {
	return issueMapKey(cloudRegion, probeName) + strconv.Itoa(int(probeAction))
}

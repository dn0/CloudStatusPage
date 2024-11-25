package charts

import (
	"context"
	"strings"

	"cspage/pkg/data"
	"cspage/pkg/mon/web/templates"
)

type issueDescription struct {
	b strings.Builder
}

func newIssueDescription() *issueDescription {
	id := &issueDescription{}
	id.b.WriteString("<ul>")
	return id
}

//nolint:varnamelen,wrapcheck,contextcheck // i(Issue) is OK in this context.
func (id *issueDescription) add(ctx context.Context, i *data.Issue) error {
	id.b.WriteString(`<li class="list-disc list-inside">`)
	if err := templates.IssueTypeSpan(i.Type).Render(ctx, &id.b); err != nil {
		return err
	}
	if err := templates.IssueSeveritySpan(i.Severity()).Render(ctx, &id.b); err != nil {
		return err
	}
	id.b.WriteString("&nbsp;")
	id.b.WriteString(i.Summary())
	id.b.WriteString("</li>")
	return nil
}

func (id *issueDescription) String() string {
	id.b.WriteString("</ul>")
	defer id.b.Reset()
	return id.b.String()
}

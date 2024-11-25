package agent

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/mon/web/views"
	"cspage/pkg/pb"
)

type ListView struct {
	dbc *db.Clients
}

func NewListView(dbc *db.Clients) *ListView {
	return &ListView{
		dbc: dbc,
	}
}

func (v *ListView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := v.handler(r)
	views.Render(w, r, c, err)
}

//nolint:wrapcheck // Error is properly logged by the caller.
func (v *ListView) handler(r *http.Request) (templ.Component, error) {
	cloud, err := data.GetCloud(chi.URLParam(r, "cloud"))
	if err != nil {
		return nil, err
	}

	agents, err := pb.GetRunningAgents(r.Context(), v.dbc.Read, cloud.Id, "")
	if err != nil {
		return nil, err
	}

	return pb.Agents(agents), nil
}

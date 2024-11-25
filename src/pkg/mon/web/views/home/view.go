package home

import (
	"net/http"

	"github.com/a-h/templ"

	"cspage/pkg/mon/web/templates"
	"cspage/pkg/mon/web/views"
)

type View struct{}

func NewView() *View {
	return &View{}
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	views.CacheControl(w, r, views.CacheMaxAgeDefault, nil)
	c, err := v.handler(r)
	views.Render(w, r, c, err)
}

//nolint:unparam // Let's keep like this for consistency
func (v *View) handler(_ *http.Request) (templ.Component, error) {
	return templates.Base(templates.NavHome, homeTempl()), nil
}

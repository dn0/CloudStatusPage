package templates

import (
	"github.com/a-h/templ"

	"cspage/pkg/data"
)

const (
	navSeparator         = "|"
	NavHome      NavName = "home"
	NavIssues    NavName = "issues"
	NavAbout     NavName = "about"
)

type NavName string

type navItem struct {
	name  NavName
	title string
	url   templ.SafeURL
}

//nolint:gochecknoglobals // This is a constant.
var navItems = newNavigation()

func newNavigation() []navItem {
	menu := []navItem{
		{
			name:  NavHome,
			title: "Home",
			url:   "/",
		},
		{
			name:  NavIssues,
			title: "Issues",
			url:   "/issues",
		},
		{
			name:  NavAbout,
			title: "About",
			url:   "/about",
		},
		{
			name: navSeparator,
		},
	}

	for _, cloud := range data.Clouds {
		menu = append(menu, navItem{
			name:  NavName(cloud.Id),
			title: cloud.Name,
			url:   templ.URL(cloud.URLPrefix()),
		})
	}

	return menu
}

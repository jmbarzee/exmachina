package service

type (
	Header struct {
		Title  string
		NavBar NavBar
	}

	NavBar struct {
		Sections []Element
	}

	Element struct {
		Title  string
		Active bool
		Link   string
	}
)

func newNavBar(activeSection string) NavBar {
	return NavBar{
		Sections: []Element{
			{
				Title:  "Dominion",
				Active: "Dominion" == activeSection,
				Link:   "/",
			},
			{
				Title:  "Devices",
				Active: "Devices" == activeSection,
				Link:   "/devices",
			},
		},
	}
}

func getTemplatesPage() []string {
	templates := []string{}

	templates = append(templates, prefixPaths("page", []string{
		"header.html",
		"navbar.html",
		"navitem.html",
		"footer.html",
	})...)

	templates = prefixPaths(staticTemplatePath, templates)
	return templates
}

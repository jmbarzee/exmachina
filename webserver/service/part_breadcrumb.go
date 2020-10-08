package service

type (
	Breadcrumb struct {
		Breadcrumbs []Element
	}
)

func getTemplatesBreadcrumb() []string {
	templates := []string{}

	templates = append(templates, prefixPaths("breadcrumb", []string{
		"breadcrumb.html",
		"breadcrumbitem.html",
	})...)

	templates = prefixPaths(staticTemplatePath, templates)
	return templates
}

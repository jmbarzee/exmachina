package service

import (
	"html/template"
	"net/http"

	"github.com/jmbarzee/dominion/identity"
)

func (s WebServer) Handle(w http.ResponseWriter, req *http.Request) {
	domains, err := s.rpcGetDomains(req.Context()) //sampleIdentities()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFiles(getTemplatesDominion()...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page := &struct {
		Header     Header
		Breadcrumb Breadcrumb
		Domains    []identity.DomainIdentity
	}{
		Header: Header{
			Title:  "Dominion",
			NavBar: newNavBar("Dominion"),
		},
		Breadcrumb: Breadcrumb{
			Breadcrumbs: []Element{
				{
					Title:  "Dominion",
					Link:   "/",
					Active: true,
				},
			},
		},
		Domains: domains,
	}
	err = t.ExecuteTemplate(w, "index.html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getTemplatesDominion() []string {
	templates := []string{}

	templates = append(templates, prefixPaths("dominion", []string{
		"index.html",
	})...)

	templates = prefixPaths(staticTemplatePath, templates)

	templates = append(templates, getTemplatesPage()...)
	templates = append(templates, getTemplatesBreadcrumb()...)

	return templates
}

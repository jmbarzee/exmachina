package service

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmbarzee/dominion/identity"
)

func (s WebServer) HandleDomain(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["domain"]

	domains, err := s.rpcGetDomains(req.Context()) //sampleIdentities()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var domain identity.DomainIdentity
	for _, d := range domains {
		if d.UUID == uuid {
			domain = d
			break
		}
	}

	t, err := template.ParseFiles(getTemplatesDomain()...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page := &struct {
		Header     Header
		Breadcrumb Breadcrumb
		Domain     identity.DomainIdentity
	}{
		Header: Header{
			Title:  "Dominion: Domain " + domain.UUID,
			NavBar: newNavBar("Dominion"),
		},
		Breadcrumb: Breadcrumb{
			Breadcrumbs: []Element{
				{
					Title: "Dominion",
					Link:  "/",
				},
				{
					Title:  domain.UUID,
					Link:   "/domain/" + domain.UUID,
					Active: true,
				},
			},
		},
		Domain: domain,
	}
	err = t.ExecuteTemplate(w, "index.html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getTemplatesDomain() []string {
	templates := []string{}

	templates = append(templates, prefixPaths("domain", []string{
		"index.html",
	})...)

	templates = prefixPaths(staticTemplatePath, templates)

	templates = append(templates, getTemplatesPage()...)
	templates = append(templates, getTemplatesBreadcrumb()...)

	return templates
}

package service

import (
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jmbarzee/dominion/ident"
)

func (s WebServer) HandleDomain(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid, err := uuid.Parse(vars["domain"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	domains, err := s.rpcGetDomains(req.Context()) //sampleIdentities()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var domain ident.DomainRecord
	for _, d := range domains {
		if d.ID == uuid {
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
		Domain     ident.DomainRecord
	}{
		Header: Header{
			Title:  "Dominion: Domain " + domain.ID.String(),
			NavBar: newNavBar("Dominion"),
		},
		Breadcrumb: Breadcrumb{
			Breadcrumbs: []Element{
				{
					Title: "Dominion",
					Link:  "/",
				},
				{
					Title:  domain.ID.String(),
					Link:   "/domain/" + domain.ID.String(),
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

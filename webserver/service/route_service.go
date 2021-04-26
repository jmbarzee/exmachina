package service

import (
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jmbarzee/dominion/ident"
)

func (s WebServer) HandleService(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid, err := uuid.Parse(vars["domain"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	stype := vars["service"]

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

	var service ident.ServiceIdentity
	for _, s := range domain.Services {
		if s.Type == stype {
			service = s
			break
		}
	}

	t, err := template.ParseFiles(getTemplatesService()...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page := &struct {
		Header     Header
		Breadcrumb Breadcrumb
		Service    ident.ServiceIdentity
	}{
		Header: Header{
			Title:  "Dominion: Domain " + domain.ID.String(),
			NavBar: newNavBar("Domains"),
		},
		Breadcrumb: Breadcrumb{
			Breadcrumbs: []Element{
				{
					Title: "Dominion",
					Link:  "/",
				},
				{
					Title: domain.ID.String(),
					Link:  "/domain/" + domain.ID.String(),
				},
				{
					Title:  service.Type,
					Link:   "/domain/" + domain.ID.String() + "/service/" + service.Type,
					Active: true,
				},
			},
		},
		Service: service,
	}
	err = t.ExecuteTemplate(w, "index.html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getTemplatesService() []string {
	templates := []string{}

	templates = append(templates, prefixPaths("service", []string{
		"index.html",
	})...)

	templates = prefixPaths(staticTemplatePath, templates)

	templates = append(templates, getTemplatesPage()...)
	templates = append(templates, getTemplatesBreadcrumb()...)

	return templates
}

package service

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmbarzee/dominion/identity"
)

func (s WebServer) HandleService(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["domain"]
	stype := vars["service"]

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

	var service identity.ServiceIdentity
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
		Service    identity.ServiceIdentity
	}{
		Header: Header{
			Title:  "Dominion: Domain " + domain.UUID,
			NavBar: newNavBar("Domains"),
		},
		Breadcrumb: Breadcrumb{
			Breadcrumbs: []Element{
				{
					Title: "Dominion",
					Link:  "/",
				},
				{
					Title: domain.UUID,
					Link:  "/domain/" + domain.UUID,
				},
				{
					Title:  service.Type,
					Link:   "/domain/" + domain.UUID + "/service/" + service.Type,
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

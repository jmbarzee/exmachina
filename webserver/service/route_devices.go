package service

import (
	"html/template"
	"net/http"
)

func (s WebServer) HandleDevices(w http.ResponseWriter, req *http.Request) {

	t, err := template.ParseFiles(getTemplatesDevices()...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page := &struct {
		Header Header
	}{
		Header: Header{
			Title:  "Dominion: Devices",
			NavBar: newNavBar("Devices"),
		},
	}
	err = t.ExecuteTemplate(w, "index.html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getTemplatesDevices() []string {
	templates := []string{}

	templates = append(templates, prefixPaths("devices", []string{
		"index.html",
	})...)

	templates = prefixPaths(staticTemplatePath, templates)

	templates = append(templates, getTemplatesPage()...)

	return templates
}

package service

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/jmbarzee/dominion/identity"
	"github.com/jmbarzee/dominion/system"
)

const (
	routeHealthCheck  string = "/healthcheck"
	routeSystemStatus string = "/systemstatus"
	routeJS           string = "/js/"
	routeCSS          string = "/css/"
	routeTMPL         string = "/tmpl/"
)

func (s WebServer) hostWebServer(ctx context.Context) {
	routineName := "hostWebServer"
	system.LogRoutinef(routineName, "Starting routine")

	LoggingFileServer := func(route string, handler http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, req *http.Request) {
				logHit(route)
				handler.ServeHTTP(w, req)
			})
	}
	mux := http.NewServeMux()

	system.Logf("%v route available", routeHealthCheck)
	mux.Handle(routeHealthCheck, http.HandlerFunc(s.HandleHealthCheck))

	system.Logf("%v route available", routeSystemStatus)
	mux.Handle(routeSystemStatus, http.HandlerFunc(s.HandleSystemStatus))

	system.Logf("%v directory available", routeJS)
	mux.Handle(routeJS, LoggingFileServer(routeJS, http.FileServer(http.Dir(path.Join(s.StaticPath, routeJS)))))

	system.Logf("%v directory available", routeCSS)
	mux.Handle(routeCSS, LoggingFileServer(routeCSS, http.FileServer(http.Dir(path.Join(s.StaticPath, routeCSS)))))

	httpServer := &http.Server{
		Addr:           ":80",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	system.Logf("Webserver listening on port 80")
	err := httpServer.ListenAndServe()
	if err != nil {
		system.Errorf("Error while serving http: %w", err)
	}
	system.LogRoutinef(routineName, "Stopping routine")
}

func (_ WebServer) HandleHealthCheck(w http.ResponseWriter, req *http.Request) {
	logHit(routeHealthCheck)
	fmt.Fprintf(w, "Healthy!")
}

func (s WebServer) HandleSystemStatus(w http.ResponseWriter, req *http.Request) {
	logHit(routeSystemStatus)
	domains, err := s.rpcGetDomains(req.Context()) //sampleIdentities()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page := &struct {
		Domains []identity.DomainIdentity
	}{
		Domains: domains,
	}
	err = s.Template.ExecuteTemplate(w, "index.html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func logHit(route string) {
	system.Logf("[Route] %v", route)
}

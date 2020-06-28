package service

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/blang/semver"
	"github.com/jmbarzee/domain/server/identity"
	"github.com/jmbarzee/domain/services"
)

type WebServer struct {
	services.Service
	Template   *template.Template
	StaticPath string
}

var htmlTemplates = []string{
	"index.html",
	"service.html",
}

func NewWebServer(port int, domainPort int, logger *log.Logger, staticPath string) (WebServer, error) {
	appendPaths := func(prefix string, paths []string) []string {
		newPaths := make([]string, len(paths))
		for i, p := range paths {
			newPaths[i] = path.Join(prefix, p)
		}
		return newPaths
	}

	logger.Printf("Loading template files...")
	template, err := template.ParseFiles(appendPaths(staticPath, appendPaths(routeTMPL, appendPaths(routeSystemStatus, htmlTemplates)))...)
	if err != nil {
		logger.Printf("Failed to load: %v", err.Error())
		return WebServer{}, err
	}
	logger.Printf("Successful!")

	logger.Printf("WebServer built!")
	return WebServer{
		Service: services.Service{
			ServiceName: "webServer",
			Port:        port,
			DomainPort:  domainPort,
			Logger:      logger,
		},
		Template:   template,
		StaticPath: staticPath,
	}, nil
}

const (
	routeHealthCheck  string = "/healthcheck"
	routeSystemStatus string = "/systemstatus"
	routeJS           string = "/js/"
	routeCSS          string = "/css/"
	routeTMPL         string = "/tmpl/"
)

func (s WebServer) LogHit(route string) {
	s.Logger.Printf("[Route] %v", route)
}

func (s WebServer) Run() {
	s.Logger.Printf("Running WebServer...")

	LoggingFileServer := func(route string, handler http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, req *http.Request) {
				s.LogHit(route)
				handler.ServeHTTP(w, req)
			})
	}
	mux := http.NewServeMux()

	s.Logger.Printf("%v route available", routeHealthCheck)
	mux.Handle(routeHealthCheck, http.HandlerFunc(s.HandleHealthCheck))

	s.Logger.Printf("%v route available", routeSystemStatus)
	mux.Handle(routeSystemStatus, http.HandlerFunc(s.HandleSystemStatus))

	s.Logger.Printf("%v directory available", routeJS)
	mux.Handle(routeJS, LoggingFileServer(routeJS, http.FileServer(http.Dir(path.Join(s.StaticPath, routeJS)))))

	s.Logger.Printf("%v directory available", routeCSS)
	mux.Handle(routeCSS, LoggingFileServer(routeCSS, http.FileServer(http.Dir(path.Join(s.StaticPath, routeCSS)))))

	httpServer := &http.Server{
		Addr:           ":" + strconv.Itoa(s.Port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.Logger.Printf("Webserver listening on %v", s.Port)
	s.Logger.Fatal(httpServer.ListenAndServe())
}

func (s WebServer) HandleHealthCheck(w http.ResponseWriter, req *http.Request) {
	s.LogHit(routeHealthCheck)
	fmt.Fprintf(w, "Healthy!")
}

func (s WebServer) HandleSystemStatus(w http.ResponseWriter, req *http.Request) {
	s.LogHit(routeSystemStatus)
	ident, idents, err := s.Dump(req.Context()) //sampleIdentities()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page := &struct {
		Identity   identity.Identity
		Identities []identity.Identity
	}{
		Identity:   ident,
		Identities: idents,
	}
	err = s.Template.ExecuteTemplate(w, "index.html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func sampleIdentities() (identity.Identity, []identity.Identity, error) {
	ident := identity.Identity{
		UUID: "UUID-ident",
		Version: semver.Version{
			Major: 1,
		},
		Services: map[string]identity.ServiceIdentity{
			"service-a": identity.ServiceIdentity{
				Port:        9001,
				LastContact: time.Now(),
			},
		},
		LastContact: time.Now(),
		IP:          net.IP{},
		Port:        9000,
	}
	return ident, []identity.Identity{ident, ident}, nil
}

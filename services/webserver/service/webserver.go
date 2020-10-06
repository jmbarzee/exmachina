package service

import (
	"context"
	"html/template"
	"path"

	"github.com/jmbarzee/dominion/service"
	"github.com/jmbarzee/dominion/service/config"
	"github.com/jmbarzee/dominion/system"
)

type WebServer struct {
	*service.Service
	Template   *template.Template
	StaticPath string
}

func NewWebServer(config config.ServiceConfig) (WebServer, error) {
	service, err := service.NewService(config)
	if err != nil {
		return WebServer{}, err
	}

	htmlTemplates := []string{
		"index.html",
		"domain.html",
		"service.html",
	}
	staticPath := "/usr/local/dominion/services/webserver/service/static"

	system.Logf("Loading template files...")
	template, err := template.ParseFiles(prefixPaths(staticPath, prefixPaths(routeTMPL, prefixPaths(routeSystemStatus, htmlTemplates)))...)
	if err != nil {
		system.Logf("Failed to load: %v", err.Error())
		return WebServer{}, err
	}
	system.Logf("Successful!")

	example := WebServer{
		Service:    service,
		Template:   template,
		StaticPath: staticPath,
	}
	return example, nil
}

func (s *WebServer) Run(ctx context.Context) error {
	system.Logf("I seek to join the Dominion\n")
	system.Logf(s.ServiceIdentity.String())
	system.Logf("The Dominion ever expands!\n")

	go s.hostWebServer(ctx)

	return s.Service.HostService(ctx)
}

func prefixPaths(prefix string, paths []string) []string {
	for i, p := range paths {
		paths[i] = path.Join(prefix, p)
	}
	return paths
}

// func sampleIdentities() (identity.Identity, []identity.Identity, error) {
// 	ident := identity.ServiceIdentity{
// 		UUID: "UUID-ident",
// 		Version: semver.Version{
// 			Major: 1,
// 		},
// 		Services: map[string]identity.ServiceIdentity{
// 			"service-a": identity.ServiceIdentity{
// 				Port:        9001,
// 				LastContact: time.Now(),
// 			},
// 		},
// 		LastContact: time.Now(),
// 		IP:          net.IP{},
// 		Port:        9000,
// 	}
// 	return ident, []identity.Identity{ident, ident}, nil
// }

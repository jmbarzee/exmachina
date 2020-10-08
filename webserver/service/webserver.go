package service

import (
	"context"
	"path"

	"github.com/jmbarzee/dominion/service"
	"github.com/jmbarzee/dominion/service/config"
	"github.com/jmbarzee/dominion/system"
)

type WebServer struct {
	*service.Service
}

const staticPath = "/usr/local/dominion/services/webserver/service/static"
const staticTemplatePath = "/usr/local/dominion/services/webserver/service/static/tmpl"

func NewWebServer(config config.ServiceConfig) (WebServer, error) {
	service, err := service.NewService(config)
	if err != nil {
		return WebServer{}, err
	}

	example := WebServer{
		Service: service,
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

package service

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmbarzee/dominion/system"
)

func (s WebServer) hostWebServer(ctx context.Context) {
	routineName := "hostWebServer"
	system.LogRoutinef(routineName, "Starting routine")

	LoggingServerFunc := func(route string, handler http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			system.Logf("[Route] %v", route)
			handler.ServeHTTP(w, req)
		}
	}
	LoggingHandlerFunc := func(route string, handlerFunc http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			system.Logf("[Route] %v", route)
			handlerFunc(w, req)
		}
	}

	r := mux.NewRouter()

	route := "/healthcheck"
	system.Logf("route available: %v", route)
	r.HandleFunc(route, LoggingHandlerFunc(route, s.HandleHealthCheck))

	route = "/"
	system.Logf("route available: %v", route)
	r.HandleFunc(route, LoggingHandlerFunc(route, s.Handle))

	route = "/domain/{domain}"
	system.Logf("route available: %v", route)
	r.HandleFunc(route, LoggingHandlerFunc(route, s.HandleDomain))

	route = "/domain/{domain}/service/{service}"
	system.Logf("route available: %v", route)
	r.HandleFunc(route, LoggingHandlerFunc(route, s.HandleService))

	route = "/devices"
	system.Logf("route available: %v", route)
	r.HandleFunc(route, LoggingHandlerFunc(route, s.HandleDevices))

	route = "/img/"
	system.Logf("route available: %v", route)
	r.PathPrefix(route).Handler(
		LoggingServerFunc(route, http.StripPrefix(route,
			http.FileServer(http.Dir(path.Join(staticPath, route))),
		)),
	)

	route = "/js/"
	system.Logf("route available: %v", route)
	r.PathPrefix(route).Handler(
		LoggingServerFunc(route, http.StripPrefix(route,
			http.FileServer(http.Dir(path.Join(staticPath, route))),
		)),
	)

	// system.Logf("route available: %v", route)
	// r.HandleFunc(routeCSS, LoggingFileServer(routeCSS, http.FileServer(http.Dir(path.Join(staticPath, routeCSS)))))

	httpServer := &http.Server{
		Addr:           ":80",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	system.Logf("Webserver listening on port 80")
	err := httpServer.ListenAndServe()
	if err != nil {
		system.Panic(fmt.Errorf("Error while serving http: %w", err))
	}
	system.LogRoutinef(routineName, "Stopping routine")
}

func (WebServer) HandleHealthCheck(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Healthy!")
}

[![Build Status](https://travis-ci.com/jmbarzee/services.svg?branch=master)](https://travis-ci.com/jmbarzee/services)
[![Go Report Card](https://goreportcard.com/badge/github.com/jmbarzee/services)](https://goreportcard.com/report/github.com/jmbarzee/services)
[![GoDoc](https://godoc.org/github.com/jmbarzee/services?status.svg)](https://godoc.org/github.com/jmbarzee/services)


# Services
Service ecosystem run by [Dominion](github.com/jmbarzee/dominion)

Service hiarchy is defined in [Dominion Config](../main/cmd/exconfig/ex.config.toml)



## WebServer
Planned: Offers a web GUI

### Routes
`/systemstatus` displays info for all services

`/healthcheck` returns `Healthy!`


## LightOrchestrator
Depenency: `WebServer`

Orchestrates all lights.


## NPLight
Depenency: `LightOrchestrator`

Peripheral: NeoPixel light string

Subscribes to light updates from the `LightOrchestrator`.


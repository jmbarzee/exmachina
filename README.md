[![License: AGPL](https://img.shields.io/badge/license-AGPL-blue.svg)](https://opensource.org/licenses/GPL-3.0/)
[![Documentation](https://godoc.org/github.com/jmbarzee/services?status.svg)](https://godoc.org/github.com/jmbarzee/services)
[![Code Quality](https://goreportcard.com/badge/github.com/jmbarzee/services)](https://goreportcard.com/report/github.com/jmbarzee/services)
[![Build Status](https://github.com/jmbarzee/services/workflows/build/badge.svg)](https://github.com/jmbarzee/services/actions)
[![Coverage](https://codecov.io/gh/jmbarzee/services/branch/master/graph/badge.svg)](https://codecov.io/gh/jmbarzee/services)


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


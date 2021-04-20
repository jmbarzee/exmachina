[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Documentation](https://godoc.org/github.com/jmbarzee/services?status.svg)](https://godoc.org/github.com/jmbarzee/services)
[![Code Quality](https://goreportcard.com/badge/github.com/jmbarzee/services)](https://goreportcard.com/report/github.com/jmbarzee/services)
[![Build Status](https://github.com/jmbarzee/services/workflows/build/badge.svg)](https://github.com/jmbarzee/services/actions)
[![Coverage](https://codecov.io/gh/jmbarzee/services/branch/master/graph/badge.svg)](https://codecov.io/gh/jmbarzee/services)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fjmbarzee%2Fservices.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fjmbarzee%2Fservices?ref=badge_shield)


# Services
Service ecosystem run by [Dominion](github.com/jmbarzee/dominion)

Service hiarchy is defined in [Dominion Config](../main/cmd/exconfig/ex.config.toml)



## WebServer
Planned: Offers a web GUI

### Routes
`/` displays dominion info

`/domain/{uuid}/` displays domain info

`/domain/{uuid}/service/{type}` displays service info

`/devices` displays device location info

`/healthcheck` returns `Healthy!`


## LightOrchestrator
Depenency: `WebServer`

Orchestrates all lights.


## NPLight
Depenency: `LightOrchestrator`

Peripheral: NeoPixel light string

Subscribes to light updates from the `LightOrchestrator`.




# Planned Work
1. Lighting Effects - Expand effects library
2. Color Picker - WebServer offers color picking for lights
3. Expose Logfiles - Services offer routes with logfiles
4. Display Logfiles - All log files can be viewed through the webserver

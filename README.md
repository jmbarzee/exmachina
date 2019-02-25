# domian
Golang Distributed System and Home Automation

[![Build Status](https://travis-ci.com/jmbarzee/domian.svg?branch=master)](https://travis-ci.com/jmbarzee/domian)
[![Go Report Card](https://goreportcard.com/badge/github.com/jmbarzee/domian)](https://goreportcard.com/report/github.com/jmbarzee/domian)
[![GoDoc](https://godoc.org/github.com/jmbarzee/domian?status.svg)](https://godoc.org/github.com/jmbarzee/domian)
[![GPL Licence](https://badges.frapsoft.com/os/gpl/gpl.svg?v=103)](https://opensource.org/licenses/GPL-3.0/)



## Purpose
This library serves on an IoT network were services (lights, speakers, thermostat, cameras, processing...) will be auto-started, auto-distributed, and (maybe) auto-scaled. domian is the back bone of this process. Once a domain starts it will start its own services, enable service discovery and distribution, and communicate about which services it still needs. 

Domains search for other domains (peers) by:
1. Identifying that they are lonely (either can't heartbeat to them or had no peers to begin with)
2. Broadcasting to the network using ZeroConf (advertizing a service with config.Title)
3. Listening for incomming RPCs (ShareIdentityList() is the heartbeat rpc)
4. Processing incomming identity lists by updating its peerMap

Domains listen for other domains (peers) by:
1. Listening for ZeroConf broadcasts (with service matching config.Title)
2. Sending ShareIdentityList RPCs to the lonely domiannaire
3. Adding the new domiannaire to its peerMap

Domains maintain contact with other domains (peers) by:
1. Checking to see if most recent contact is too old
2. Checking and possible establishing a new connection with the domiannaire (peer)
3. Sending ShareIdentityList RPCs to the domiannaire (peer)
4. Processing the replied identity list by updating its peerMap

Domains processing identity lists by:
1. Adding new identities as peers without connections
2. Updating current peers information, like IP, Port, LastContact, and ServiceList



## Utilized Libraries

`github.com/blang/semver`

`google.golang.org/grpc`

`github.com/grandcat/zeroconf`

`github.com/BurntSushi/toml`



## Planned Development

1. Service Sharing - strategies for service start and service dependency evaluation 
2. Connection encryption - encrypt RPCs
3. Identity verification - sign communication with preestablished keypairs




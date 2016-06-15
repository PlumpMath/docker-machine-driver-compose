Overview
========

Docker machine driver for brooklyn

Getting Started
===============

Installing GO
-------------

Follow the instructions mention in getting started on [Installing Go](https://golang.org/doc/)


Build The Driver
----------------
- go get ./...
- go build stash.fsc.atos-services.net/scm/cet/bdmd.git
- go run stash.fsc.atos-services.net/scm/cet/bdmd.git

Development Environment
-----------------------
`a588232@MC0WBVEC ~/Jitendra/Workspace/go.ws/bin
$ go build ../src/stash.fsc.atos-services.net/scm/cet/bdmd.git/docker-machine
-driver-brooklyn.go`

`a588232@MC0WBVEC ~/Jitendra/Workspace/go.ws/bin
$ docker-machine create --driver brooklyn  --brooklyn-base-url https://test.c
ompose.canopy-cloud.com --brooklyn-user compose.test@canopy-cloud.com --brook
lyn-password password --brooklyn-target-location "AWS Frankfurt"  machinename`
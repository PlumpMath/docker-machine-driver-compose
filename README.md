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

Build Latest Driver
`$ go build ../src/stash.fsc.atos-services.net/scm/cet/bdmd.git/docker-machine-driver-brooklyn.go`

Create DockerHost Without Swarm Manager
`$ docker-machine create --driver brooklyn  \
    --brooklyn-base-url https://test.compose.canopy-cloud.com \
    --brooklyn-user compose.test@canopy-cloud.com 
    --brooklyn-password password --brooklyn-target-location "AWS Frankfurt" machinename`
    
Create Docker Swarm Manager
`$ docker-machine create --driver brooklyn  \
    --brooklyn-base-url https://test.compose.canopy-cloud.com \
    --brooklyn-user compose.test@canopy-cloud.com \
    --brooklyn-password password --brooklyn-target-location "AWS Frankfurt" \ 
    --swarm --swarm-master --swarm-discovery token://SWARM_CLUSTER_TOKEN \    
    swarm-manager`
    
Create Docker Host With Registering Swarm Manager
`$ docker-machine create --driver brooklyn  \
    --brooklyn-base-url https://test.compose.canopy-cloud.com \
    --brooklyn-user compose.test@canopy-cloud.com \
    --brooklyn-password password --brooklyn-target-location "AWS Frankfurt" \ 
    --swarm --swarm-discovery token://SWARM_CLUSTER_TOKEN \    
    node-01`    
    

Test Newly created Dockerhost
`$ docker --tlsverify --tlscacert=/c/Users/A588232/.docker/machine/certs/ca.pem \ 
    --tlscert=/c/Users/A588232/.docker/machine/certs/cert.pem \
    --tlskey=/c/Users/A588232/.docker/machine/certs/key.pem \
    -H=ec2-52-28-2-68.eu-central-1.compute.amazonaws.com:2376 version`
    
Docker Swarm Manager Info
`$ docker --tlsverify --tlscacert=/c/Users/A588232/.docker/machine/certs/ca.pem \
    --tlscert=/c/Users/A588232/.docker/machine/certs/cert.pem \
    --tlskey=/c/Users/A588232/.docker/machine/certs/key.pem \
    -H=ec2-52-59-20-162.eu-central-1.compute.amazonaws.com:3376 info`    
    
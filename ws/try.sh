#!/bin/bash

export CTR_REGISTRY=docker.dev.ws:5000
export CTR_TAG=osmv9.3
#make docker-push
make docker-push-init-osm-controller
make docker-push-osm-controller
make docker-push-init
make docker-push-osm-injector

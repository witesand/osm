#!/bin/bash

export CTR_REGISTRY=docker.dev.ws:5000
export CTR_TAG=osmv9.2
make docker-push-init-osm-controller
make docker-push-osm-controller
make docker-push-init

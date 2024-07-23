#!/bin/bash

export $(grep -v '^#' .env | xargs)

envsubst < config/prod.docker.json > prod.docker.tmp.json && mv prod.docker.tmp.json config/prod.docker.json
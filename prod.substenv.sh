#!/bin/bash

export $(grep -v '^#' .env | xargs)

envsubst < configs/prod.docker.json > prod.docker.tmp.json && mv prod.docker.tmp.json configs/prod.docker.json
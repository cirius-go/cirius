#!/bin/bash

source ./.env
source ./.env.local

./scripts/initdb.sh

docker compose build

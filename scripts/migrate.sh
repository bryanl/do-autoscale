#!/usr/bin/env bash


DB_URL=postgres://${AUTOSCALE_DB_USER}:${AUTOSCALE_DB_PASSWORD}@${AUTOSCALE_DB_HOST}:${AUTOSCALE_DB_PORT}/${AUTOSCALE_DB_NAME}?sslmode=disable

migrate -url $DB_URL -path src/autoscale/db/migrations $@
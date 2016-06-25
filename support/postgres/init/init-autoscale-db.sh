#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE USER autoscale;
    CREATE DATABASE autoscale;
    GRANT ALL PRIVILEGES ON DATABASE autoscale TO autoscale;
EOSQL
#!/bin/bash

docker exec -t -e PGPASSWORD=secret store-management_db_1 pg_dump -U root  -h localhost -p 5432 -d store_management -f backup/db.sql


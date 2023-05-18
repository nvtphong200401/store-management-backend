#!/bin/bash
docker exec -t -e PGPASSWORD=secret store-management_db_1 psql -U root  -h localhost -p 5432 -d store_management -c "DROP SCHEMA public CASCADE"
docker exec -t -e PGPASSWORD=secret store-management_db_1 psql -U root  -h localhost -p 5432 -d store_management -c "CREATE SCHEMA public"
docker exec -t -e PGPASSWORD=secret store-management_db_1 psql -U root  -h localhost -p 5432 -d store_management -f backup/db.sql
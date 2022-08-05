
# Run make compose-db, then wait for a bit and run make populate-db before running make run
# 
# Alternatively to make run, make compose will run the server with docker

SCYLLA_HOST=scylla-node1
SCYLLA_DIR=$(shell pwd)/database/scylla
SCYLLA_NETWORK=scylla_database

POSTGRES_HOST=postgres
POSTGRES_DIR=$(shell pwd)/database/postgres
POSTGRES_NETWORK=postgres_database

install:
	go mod download && go mod verify

run:
	go run $(shell pwd)/src/server.go

compose:
	docker-compose up --build

kill:
	docker-compose down

compose-scylla:
	docker-compose -f ${SCYLLA_DIR}/docker-compose.yml up

populate-scylla:
	docker run -it --network ${SCYLLA_NETWORK} --env CQLSH_HOST=${SCYLLA_HOST} -v "${SCYLLA_DIR}:/host-data/" --rm cassandra cqlsh -f /host-data/create-scylla.cql

kill-scylla:
	docker-compose -f ${SCYLLA_DIR}/docker-compose.yml down

compose-postgres:
	docker-compose -f ${POSTGRES_DIR}/docker-compose.yml up

create-postgres:
	docker run -it --network ${POSTGRES_NETWORK} --env PGPASSWORD=password -v "${POSTGRES_DIR}:/host-data/" postgres psql -h ${POSTGRES_HOST} -U postgres -f /host-data/create-db.sql

populate-postgres:
	docker run -it --network ${POSTGRES_NETWORK} --env PGPASSWORD=password -v "${POSTGRES_DIR}:/host-data/" postgres psql -h ${POSTGRES_HOST} -U postgres -f /host-data/populate-db.sql

kill-postgres:
	docker-compose -f ${POSTGRES_DIR}/docker-compose.yml down

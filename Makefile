
DB_HOST=scylla-node1
DB_COMPOSE=$(shell pwd)/database/scylla-compose.yml
DB_NETWORK=database_database

install:
	go mod download && go mod verify

run:
	docker-compose up --build

kill:
	docker-compose down

run-db:
	docker-compose -f ${DB_COMPOSE} up

populate-db:
	docker run -it --network ${DB_NETWORK} --env CQLSH_HOST=${DB_HOST} -v "$(shell pwd)/database:/host-data/" --rm cassandra cqlsh -f /host-data/create-scylla.cql

kill-db:
	docker-compose -f ${DB_COMPOSE} down

kill-all:
	docker-compose down
	docker-compose -f ${DB_COMPOSE} down

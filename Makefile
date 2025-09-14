include ./.env
MIGRATIONPATH=db/migrations
DBURL=postgres://$(DBUSER):$(DBPASS)@$(DBHOST):$(DBPORT)/$(DBNAME)?sslmode=disable

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONPATH) -seq create_$(NAME)_table

migrate-up:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) up

migrate-down:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) down

print-db-url:
	echo $(DBURL)
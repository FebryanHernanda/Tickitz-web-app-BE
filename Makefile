include ./.env
MIGRATIONPATH=db/migrations
DBURL=postgres://$(DBUSER):$(DBPASSWORD)@$(DBHOST):$(DBPORT)/$(DBNAME)?sslmode=disable

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONPATH) -seq create_$(NAME)_table

migrate-up:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) up

migrate-down:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) down

migrate-seed:
	psql $(DBURL) -f $(SEEDPATH)/$(NAME).sql

migrate-status:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) version

migrate-force:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) force $(v)
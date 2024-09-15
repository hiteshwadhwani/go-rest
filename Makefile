db-start:
	@mkdir -p testdata/postgres
	docker run --rm --name postgres-go-rest -v $(shell PWD)/testdata/:/testdata \
	-v $(shell PWD)/testdata/postgres:/var/lib/postgresql/data \
	-e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=go-rest -d -p 5432:5432 postgres

db-stop:
	docker stop postgres-go-rest

setup:
	docker run --name sns-db -e POSTGRES_PASSWORD=postgres -p 5434:5432 -d postgres

create-db:
	docker exec -it sns-db psql -U postgres -c "CREATE DATABASE events;"

drop-db:
	docker exec -it sns-db psql -U postgres -c "DROP DATABASE events;"

teardown:
	docker stop sns-db
	docker rm sns-db

connect-db:
	psql postgresql://postgres:postgres@localhost:5434/events
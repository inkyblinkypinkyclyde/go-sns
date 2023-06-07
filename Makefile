setup-dc:
	docker-compose down
	docker-compose up -d

connect-db:
	psql postgresql://postgres:postgres@localhost:5434/events

create-volumes:
	docker volume create --name=go_sns_volumes

remove-volumes:
	docker rm go_sns_volumes
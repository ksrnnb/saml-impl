run: keycloak.run
	go run ./

keycloak.run:
	docker-compose up -d

down: keycloak.down

keycloak.down:
	docker-compose down

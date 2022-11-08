run: keycloak.run
	go run ./

keycloak.run:
	docker run -d -p 8080:8080 -e KEYCLOAK_ADMIN=admin -e KEYCLOAK_ADMIN_PASSWORD=admin quay.io/keycloak/keycloak:20.0.0 start-dev

down: keycloak.down

keycloak.down:
	docker ps | grep keycloak | cut -d ' ' -f 1 | xargs docker stop

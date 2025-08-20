integration:
		@docker-compose stop zookeeper kafka postgres || true
		@docker-compose -f ./internal/tests/integration/docker-compose.yaml up -d
		@sleep 20
		@status=0; \
    	go test -race -v -timeout 200s ./internal/tests/integration/consumer_test.go -tags=integration || status=$$?; \
    	go test -race -v -timeout 200s ./internal/tests/integration/http_test.go -tags=integration || status=$$?; \
    	docker-compose -f ./internal/tests/integration/docker-compose.yaml down --volumes --remove-orphans; \
    	docker-compose start zookeeper kafka postgres || true; \
    	exit $$status

migrate-up:
	@docker compose run --rm app sh -c './goose -dir migrations postgres "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSL" up'

migrate-down:
	@docker compose run --rm app sh -c './goose -dir migrations postgres "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSL" down'

run:
	@docker compose up -d zookeeper kafka postgres
	@$(MAKE) migrate-up
	@docker compose up app

down:
	@docker compose down --volumes --remove-orphans

logs:
	@docker compose logs -it -f app

demo:
	@docker build -t kafka-producer ./demo-producer
	@docker run --rm -it --network order-service_dev_network kafka-producer


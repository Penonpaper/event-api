.PHONY: test-integration test-integration-down

test-integration:
	docker compose up -d --wait postgres_test
	docker compose --profile integration run --rm migrator_test
	go test -tags=integration ./...
	docker compose stop postgres_test

test-integration-down:
	docker compose stop postgres_test
	docker compose rm -f postgres_test
	docker compose rm -f book-api_postgres_test_data
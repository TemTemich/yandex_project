

start_with_logs:
	@docker compose -f docker-compose.yaml --env-file config/.env up



down:
	@docker compose -f docker-compose.yaml --env-file config/.env down --remove-orphans


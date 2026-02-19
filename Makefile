#######################
# Docker Compose
#######################
.PHONY: dc-up
dc-up:
	docker compose -f docker-compose.yml up -d

.PHONY: dc-o11y-up
dc-o11y-up:
	docker compose -f docker-compose.o11y.yml up -d

.PHONY: dc-all-up
dc-all-up:
	docker compose -f docker-compose.yml -f docker-compose.o11y.yml up -d

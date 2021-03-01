up:
	docker-compose -f docker-compose.yml up -d --build

down:
	docker-compose -f docker-compose.yml down

test:
	make down
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
	docker-compose -f docker-compose.test.yml down

test-db-up:
	docker-compose -f docker-compose.test.yml up --build database

test-db-down:
	docker-compose -f docker-compose.test.yml down --volumes database
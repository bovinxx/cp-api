SERVER_IMAGE_NAME = pudge
SERVER_PORT = 8000
RABBITMQ_IMAGE_NAME = rabbitmq:latest
RABBITMQ_PORT = 5672

BUILD = docker-compose build
RUN = docker-compose up
BUILD_CODE_PROCESSOR_BASIC_IMAGE = docker build -f ./processor/basic_image/Dockerfile -t basic_image .
GENERATE_DOCS = swag init --output api/docs -g api/cmd/app/main.go
SET_TMP_DIR = python util/make_tmp_dir.py
MIGRATE_UP = goose -dir postgres/migrations postgres "postgresql://user:password@127.0.0.1:5432/db?sslmode=disable" up
MIGRATE_DOWN = goose -dir postgres/migrations postgres "postgresql://user:password@127.0.0.1:5432/db?sslmode=disable" down

basic_image:
	$(BUILD_CODE_PROCESSOR_BASIC_IMAGE)

migrate_up:
	$(MIGRATE_UP)

migrate_down:
	$(MIGRATE_DOWN)

build:
	$(BUILD) 

run:
	$(RUN)

generate_docs:
	$(GENERATE_DOCS)

tmp_dir:
	$(SET_TMP_DIR)

start:
	$(BUILD)
	$(RUN)
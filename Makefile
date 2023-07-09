# .PHONYを書くとダミーターゲットを作ることができま
# ダミーターゲットを使うと、そのダミーターゲットをビルドするためのコマンドを実行することができます。
# 実際にはターゲットはダミーなので、コマンドのみが実行されます。
# つまりMakefileでコマンドだけを実行したい場合は、この.PHONYでダミーターゲットを作れば可能だということになります。
# .PHONYを記載しないとコマンドと同じ名前のファイルがある場合衝突してしまいコマンドが実行できない
.PHONY: help build build-local up down logs ps test
.DEFAULT_GOAL := help

DOCKER_TAG := latest
build:
	docker compose build --no-cache

# Build docker image to deploy
up:
	docker compose up

down:
	docker compose down

dry-migrate: ## Try migration
	mysqldef -u todo -p todo -h 127.0.0.1 -P 33306 todo --dry-run < ./_tools/mysql/schema.sql

migrate:  ## Execute migration
	mysqldef -u todo -p todo -h 127.0.0.1 -P 33306 todo < ./_tools/mysql/schema.sql

generate:
	go generate ./...
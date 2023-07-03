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
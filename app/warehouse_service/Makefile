GO_CMD_MAIN = cmd/*.go

run:
	go run $(GO_CMD_MAIN) server

generate-sqlc:
	sqlc generate
	go run gen.go

migrate:
	echo \# make migrate name="$(name)"
	go run $(GO_CMD_MAIN) migrate create $(name)

migrate-up:
	go run $(GO_CMD_MAIN) migrate up

migrate-down-1:
	go run $(GO_CMD_MAIN) migrate down 1

test:
	go test ./...  -count=1 -v -cover

.PHONY = run generate-sqlc migrate migrate-up test

proto:
	protoc -I./api --go_out=./api --go_opt=paths=source_relative \
        --go-grpc_out=./api --go-grpc_opt=paths=source_relative \
        api/api.proto
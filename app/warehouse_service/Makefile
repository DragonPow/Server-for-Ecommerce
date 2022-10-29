GO_CMD_MAIN = cmd/*.go

run:
	go run $(GO_CMD_MAIN) server

update:
	cd ../.. && make update

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
	cd ../.. && protoc \
		-I . \
		-I vendor \
		-I third_party/envoyproxy \
		-I third_party/gogoprotobuf \
		-I third_party/googleapis \
		--gofast_out=plugins=grpc,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:. \
		--validate_out=lang=go:. \
		--grpc-gateway_out=allow_repeated_fields_in_body=true:. \
		--ecode_out=. \
        app/warehouse_service/api/api.proto
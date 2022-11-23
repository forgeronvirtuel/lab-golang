.PHONY: build

build: clean
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		./api/route.proto
	mkdir -p ./build
	go build -o ./build/server ./cmd/server.go
	@echo "---- BUILD OK ----"

clean:
	rm -f api/*.pb.go
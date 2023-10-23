# CLIENT
build-client:
	@cd client && go build -o ./bin/vincent

run-client: build-client
	@./client/bin/vincent

test-client:
	@go test -v ./server/...

# SERVER
build-server:
	@cd server && go build -o ./bin/vincent

run-server: build-server
	@./server/bin/gobank

test-server:
	@go test -v ./client/...


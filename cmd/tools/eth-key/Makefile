
.PHONY:build
build:
	go build -o eth-key main.go

exec:
	./eth-key -keydir ./tmp -keyfile ./keys/key-sample.json -password password

.PHONY:run
run:
	go run main.go -keydir ./tmp -keyfile ./keys/key-sample.json -password password

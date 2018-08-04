goget:
	go get -u -d -v ./...

bld:
	go build -o detect ./cmd/detectinput/main.go

# Testing directory path
tpath:=./dat/testing
# List of packages to test; -e [package] excludes from testing
tlist:=`go list ./... | grep -v -e authserver -e cmd`

test:
	go test $(tlist) -coverprofile $(tpath)/coverage.profile
	go tool cover -html=$(tpath)/coverage.profile -o $(tpath)/coverage.html

server:
	go build -o server ./cmd/server/main.go
	mv ./server ./bin

server-run:
	make server
	./bin/server

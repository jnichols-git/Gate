tpath:=./dat/testing
tlist:=`go list ./...`

make test:
	go test $(tlist) -coverprofile $(tpath)/coverage.profile
	go tool cover -html=$(tpath)/coverage.profile -o $(tpath)/coverage.html

tpath:=./dat/testing # Output for test results
tlist:=`go list ./...` # List of files to test
# ssl cert
sslC:=US
sslST:=Colorado
sslL:=Boulder
sslO:=jakenichols2719
ssl

make test:
	go test $(tlist) -coverprofile $(tpath)/coverage.profile
	go tool cover -html=$(tpath)/coverage.profile -o $(tpath)/coverage.html

make cert:
	openssl

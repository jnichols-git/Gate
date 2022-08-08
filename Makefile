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

docker-image:
	docker build -t jakenichols2719/gate .

kube-apply:
	kubectl apply -f kubernetes/

kube-expose:
	kubectl expose deployment gate --type=LoadBalancer --name=gate-service --port=443

kube-restart:
	kubectl rollout restart deployment/gate

kube-stop:
	kubectl delete deployment gate
	kubectl delete service gate-service

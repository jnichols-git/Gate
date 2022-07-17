package main

import (
	"auth/pkg/authserver"
	"errors"
	"time"
)

var secret []byte = []byte("test secret")

func main() {
	srv, err := authserver.ServerFromConfig("./dat/config/config.yaml")
	if srv == nil {
		panic(errors.New("Failed to create AuthServer"))
	}
	if err != nil {
		panic(err)
	}
	srv.Start()
	time.Sleep(time.Minute * 1)
	srv.Stop()
	time.Sleep(time.Second * 1)
}

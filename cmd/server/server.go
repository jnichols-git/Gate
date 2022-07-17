package main

import (
	"auth/pkg/authserver"
	"errors"
	"fmt"
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
	time.Sleep(time.Second * 1)
	fmt.Println("Stopping server")
	srv.Stop()
}

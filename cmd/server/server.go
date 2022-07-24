package main

import (
	"auth/pkg/server"
	"fmt"
	"time"
)

var secret []byte = []byte("test secret")

func main() {
	// Create and start server
	cfg := server.NewConfig()
	err := cfg.ReadConfig("./dat/config/config.yaml")
	if err != nil {
		panic(err)
	}
	srv := server.NewServer(cfg)
	srv.Start()
	// Sleep for 5 minutes for testing time
	time.Sleep(time.Minute * 30)
	fmt.Println("Stopping server")
	srv.Stop()
	fmt.Println("Done!")
}

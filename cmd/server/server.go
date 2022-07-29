package main

import (
	"fmt"
	"time"

	"github.com/jakenichols2719/gate/pkg/server"
)

var secret []byte = []byte("test secret")

func main() {
	// Create and start server
	cfg := server.NewConfig()
	err := cfg.ReadConfig("./dat/config/config.yml")
	if err != nil {
		panic(err)
	}
	srv := server.NewServer(cfg)
	err = srv.Start()
	if err != nil {
		panic(err)
	}
	// Sleep for 5 minutes for testing time
	time.Sleep(time.Minute * 30)
	fmt.Println("Stopping server")
	srv.Stop()
	fmt.Println("Done!")
}

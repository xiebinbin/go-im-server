package main

import (
	"imsdk/internal/client"
	"imsdk/pkg/server"
	"os"
)

func main() {
	err := os.Setenv("RUN_MODULE", "client")
	if err != nil {
		return
	}
	err = os.Setenv("SOCKET_HOST", "127.0.0.1")
	if err != nil {
		return
	}
	server.Mode = "client"
	server.Run(client.GetEngine)
}

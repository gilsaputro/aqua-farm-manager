package main

import (
	"aqua-farm-manager/cmd/aqua-farm-manager/server"
	"os"
)

func main() {
	os.Exit(server.Run())
}

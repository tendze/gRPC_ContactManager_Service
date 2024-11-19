package main

import (
	"fmt"
	"gRPC_ContactManagement_Service/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: LOAD CONFIG
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO: LOGGER

}

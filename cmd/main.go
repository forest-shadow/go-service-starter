package main

import (
	"fmt"
	"repo-starter/internal/config"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("Error while getting config: %s", err)
	}
	fmt.Printf("%+v", cfg)
}

package main

import (
	"fmt"
	"log"

	"github.com/CTSDM/gator-go/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reding file: %v", err)
	}
	fmt.Printf("Read config: %v\n", cfg)

	cfg.SetUser("victor")

	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Reading modified config: %v\n", cfg)
}

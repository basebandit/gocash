package main

import (
	"fmt"
	"log"
	"os"

	"github.com/basebandit/gocash/pkg/api"
	"github.com/basebandit/gocash/pkg/cash"
	"github.com/basebandit/gocash/pkg/config"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	configFile := fmt.Sprintf("%s/config.json", wd)

	cfg, err := config.LoadConfig(configFile)

	if err != nil {
		log.Fatal(err)
	}

	if cfg.ApiKey == "" {
		log.Fatal("Missing API Key")
	}

	url := fmt.Sprintf(cfg.Api, cfg.ApiKey)

	body, err := api.Fetch(url)
	if err != nil {
		log.Fatal(err)
	}
	if rates := cash.ParseRates(body); rates != nil {
		fmt.Println(rates)
	}
}

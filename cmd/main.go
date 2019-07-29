package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/basebandit/gocash/pkg/api"
	"github.com/basebandit/gocash/pkg/cash"
	"github.com/basebandit/gocash/pkg/config"
	"github.com/fatih/color"
)

var configDir = flag.String("config", "", "Configuration directory path")
var red = color.New(color.FgRed)
var boldRed = red.Add(color.Bold)

// var from = flag.Bool("from", false, "Enable debug output")
// var to = flag.Bool("to", false, "Enable developer mode (generates self-signed certificates for all hostnames)")

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		help()
		os.Exit(1)
	}

	amount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	from := args[1]

	to := args[2]

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	parentDir := filepath.Dir(wd)
	configFile := fmt.Sprintf("%s/config.json", parentDir)

	cfg, err := config.LoadConfig(configFile)

	if err != nil {
		log.Fatal(err)
	}

	if cfg.ApiKey == "" {
		boldRed.Println("\nFixer.io API key not found!")
		color.Cyan("Get it here for free: https://fixer.io/signup/free")
		color.Cyan("Then run `cash --key [key]` to save it\n")
		os.Exit(1)
	}

	url := fmt.Sprintf(cfg.Api, cfg.ApiKey)

	body, err := api.Fetch(url)
	if err != nil {
		log.Fatal(err)
	}

	money := new(api.Money)

	if rates := cash.ParseRates(body); rates != nil {
		money.Rates = rates
	}

	if base := cash.ParseBase(body); base != "" {
		money.Base = base
	}

	//convert
	amt, err := money.Convert(amount, from, to)
	color.Green(strconv.FormatFloat(amt, 'f', 6, 64))
}

func help() {
	help := `
	Usage
		$ cash <amount> <from> <to>
		$ cash <options>
	Options
	  	--config -c   config file
		--save -s 			Save default currencies
		--purge -p 			Purge cached API response to get the latest data
	Examples
		$ cash --key [key]
		$ cash 10 usd eur pln
		$ cash --save usd aud 
	`
	fmt.Println(help)
}

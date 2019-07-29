package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/basebandit/gocash/pkg/api"
	"github.com/basebandit/gocash/pkg/config"
	"github.com/basebandit/gocash/pkg/currency"
	"github.com/fatih/color"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
)

var configDir = flag.String("config", "", "Configuration directory path")
var red = color.New(color.FgRed)
var boldRed = red.Add(color.Bold)

func main() {
	args := os.Args[1:]

	if len(args) < 3 {
		help()
		os.Exit(1)
	}

	amount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	from := strings.ToUpper(args[1])

	toItems := args[2:]

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	parentDir := filepath.Dir(wd)
	configFile := fmt.Sprintf("%s/config.json", parentDir)

	cfg, err := config.LoadConfig(configFile)

	if err != nil {
		boldRed.Println(err.Error())
		os.Exit(1)
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
		boldRed.Println(err.Error())
		os.Exit(1)
	}

	money := new(api.Money)

	if rates := currency.ParseRates(body); rates != nil {
		money.Rates = rates
	}

	if base := currency.ParseBase(body); base != "" {
		money.Base = base
	}

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	//start loader
	w := wow.New(os.Stdout, spin.Get(spin.Dots), green("  Converting ..."))
	w.Start()
	time.Sleep(2 * time.Second)

	currencyFile := fmt.Sprintf("%s/currencies.json", parentDir)
	currencies, err := config.GetCurrencies(currencyFile)

	if err != nil {
		boldRed.Println(err.Error())
		os.Exit(1)
	}

	//Validate from currency
	if _, ok := currencies[from]; !ok {
		w.PersistWith(spin.Spinner{Frames: []string{"👍  "}}, yellow(fmt.Sprintf("The %s currency not found\n", from)))
	}

	for _, t := range toItems {
		to := strings.ToUpper(t)
		//Validate to currency
		if _, ok := currencies[to]; !ok {
			w.PersistWith(spin.Spinner{Frames: []string{"👎  "}}, yellow(fmt.Sprintf("The %s currency not found\n", to)))
		}
		//convert
		amt, err := money.Convert(amount, from, to)

		if err != nil {
			boldRed.Println(err.Error())
			os.Exit(1)
		}
		w.PersistWith(spin.Spinner{Frames: []string{"👍  "}}, fmt.Sprintf("%s (%s) %s\n", green(strconv.FormatFloat(amt, 'f', 2, 64)), to, currencies[to]))
	}

	c := color.New(color.FgHiBlack).Add(color.Bold).Add(color.Underline)
	c.Println(fmt.Sprintf("Conversion of %s %.f \n", from, amount))
}

func help() {
	help := `
	Usage
		$ cash <amount> <from> <to>
		$ cash <options>
	Options
	  	--config -c   config file
	Examples
		$ cash --key [key]
		$ cash 10 usd eur 
	`
	fmt.Println(help)
}

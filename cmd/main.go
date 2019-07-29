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
	"github.com/basebandit/gocash/pkg/cash"
	"github.com/basebandit/gocash/pkg/config"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/fatih/color"
)

var configDir = flag.String("config", "", "Configuration directory path")
var red = color.New(color.FgRed)
var boldRed = red.Add(color.Bold)

// var from = flag.Bool("from", false, "Enable debug output")
// var to = flag.Bool("to", false, "Enable developer mode (generates self-signed certificates for all hostnames)")

func main() {
	args := os.Args[1:]

	if (len(args) < 3) || (len(args) > 3) {
		help()
		os.Exit(1)
	}

	amount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}


	from := strings.ToUpper(args[1])

	to := strings.ToUpper(args[2])

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

	green := color.New(color.FgGreen).SprintFunc()

	//start loader 
	w := wow.New(os.Stdout, spin.Get(spin.Dots), green("  Converting ..."))
	w.Start()
	time.Sleep(2 * time.Second)
	
	//convert
	amt, err := money.Convert(amount, from, to)

	w.PersistWith(spin.Spinner{Frames: []string{"👍  "}},fmt.Sprintf("%s\n",green(strconv.FormatFloat(amt, 'f', 6, 64))))

	c := color.New(color.FgMagenta).Add(color.Bold).Add(color.Underline)
	c.Println(fmt.Sprintf("Conversion of %s %.2f to %s\n",from,amount,to))
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

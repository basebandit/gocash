package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/basebandit/gocash/pkg/config"
	"github.com/basebandit/gocash/pkg/currency"
	"github.com/fatih/color"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/mitchellh/go-homedir"
)

//red color
var red = color.New(color.FgRed)

//boldRed color for fatal error
var boldRed = red.Add(color.Bold)

//greeen color for success messages
var green = color.New(color.FgGreen).SprintFunc()

//yellow color for warnings
var yellow = color.New(color.FgYellow).SprintFunc()

//loader to indicate progress
var w = wow.New(os.Stdout, spin.Get(spin.Dots), green(" Converting ..."))

//RunConverter starts the currency conversion
func RunConverter() {

	args := os.Args[1:]

	if len(args) < 3 {
		help()
		os.Exit(1)
	}

	amount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		boldRed.Println(err.Error())
		os.Exit(1)
	}

	from := strings.ToUpper(args[1])

	toItems := args[2:]

	configDir := createConfigDir(".gocash")

	configFileURL := "https://github.com/basebandit/gocash/blob/master/config.json"

	currencyFileURL := "https://github.com/basebandit/gocash/blob/master/currencies.json"

	configFile := fmt.Sprintf("%s/config.json", configDir)

	cfg, err := config.LoadConfig(configFile)

	if err != nil {
		boldRed.Printf("Could not find your config file.Copy %s inside %s\n", yellow(configFileURL), yellow(configDir))
		os.Exit(1)
	}

	currencyFile := fmt.Sprintf("%s/currencies.json", configDir)
	currencies, err := config.GetCurrencies(currencyFile)

	if err != nil {
		boldRed.Printf("Could not find currencies list.Copy %s inside %s\n", yellow(currencyFileURL), yellow(configDir))
		os.Exit(1)
	}

	if cfg.ApiKey == "" {
		boldRed.Println("\nFixer.io API key not found!")
		color.Cyan("Get it here for free: https://fixer.io/signup/free")
		color.Cyan("Then save it inside `.gocash/config.json` in your home directory\n")
		os.Exit(1)
	}

	money, err := getCurrencyRates(cfg)

	if err != nil {
		boldRed.Println(err.Error())
		os.Exit(1)
	}
	//start loader
	w.Start()
	time.Sleep(2 * time.Second)

	convertCurrency(currencies, amount, from, toItems, money)

	c := color.New(color.FgHiBlack).Add(color.Bold).Add(color.Underline)
	c.Println(fmt.Sprintf("Conversion of %s %.f \n", from, amount))
}

//help prints the help menu
func help() {
	help := `
	Usage
		$ cash <amount> <from> <to>

	Examples
		$ cash 10 usd eur 
		$ cash 100 eur usd kes tzs aud zwl vnd
	`
	fmt.Println(help)
}

//joinPath joins two relative path to form a full path
func joinPath(basePath string, relPath string) string {
	var retPath string
	if filepath.IsAbs(relPath) {
		retPath = relPath
	} else {
		retPath = filepath.Join(basePath, relPath)
	}
	return retPath
}

//mkDir creates .gocash directory in the user's home directory
func createConfigDir(configDir string) string {
	var configDirPath string

	//Get user home directory
	homeDir, err := homedir.Dir()

	if err != nil {
		boldRed.Println(err.Error())
		os.Exit(1)
	}

	configDirPath = joinPath(homeDir, configDir)

	//Check if config Dir already exists
	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
		//Lets create the directory if it does not exist
		err := os.MkdirAll(configDirPath, os.FileMode(0700))
		if err != nil {
			boldRed.Println(err.Error())
			os.Exit(1)
		}
	}
	return configDirPath
}

//getCurrencyRates retrieves the currency rates from the json response
func getCurrencyRates(cfg *config.Config) (*currency.Currency, error) {
	url := fmt.Sprintf(cfg.Api, cfg.ApiKey)

	body, err := currency.Fetch(url)
	if err != nil {
		return nil, err
	}

	money := new(currency.Currency)

	if rates := currency.ParseRates(body); rates != nil {
		money.Rates = rates
	}

	if base := currency.ParseBase(body); base != "" {
		money.Base = base
	}
	return money, nil
}

func convertCurrency(currencies map[string]interface{}, amount float64, from string, toItems []string, money *currency.Currency) {
	//Validate from currency
	if _, ok := currencies[from]; !ok {
		w.PersistWith(spin.Spinner{Frames: []string{"\xE2\x9C\x97  "}}, yellow(fmt.Sprintf("The %s currency not found\n", from)))
		os.Exit(1)
	}

	for _, t := range toItems {
		to := strings.ToUpper(t)
		//Validate to currency
		if _, ok := currencies[to]; !ok {
			w.PersistWith(spin.Spinner{Frames: []string{"\xE2\x9C\x97  "}}, yellow(fmt.Sprintf("The %s currency not found\n", to)))
			os.Exit(1)
		}
		//convert
		amt, err := money.Convert(amount, from, to)

		if err != nil {
			boldRed.Println(err.Error())
			os.Exit(1)
		}
		w.PersistWith(spin.Spinner{Frames: []string{"\xE2\x9C\x93  "}}, fmt.Sprintf("%s (%s) %s\n", green(strconv.FormatFloat(amt, 'f', 2, 64)), to, currencies[to]))
	}
}

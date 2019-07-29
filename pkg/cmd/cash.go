package cmd

import (
	"flag"
	"fmt"
	"io"
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

func RunConverter() {
	red := color.New(color.FgRed)
	boldRed := red.Add(color.Bold)
	args := os.Args[1:]

	if len(args) < 3 {
		help()
		os.Exit(1)
	}

	flag.Parse() // add this line

	amount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		boldRed.Println(err.Error())
		os.Exit(1)
	}

	from := strings.ToUpper(args[1])

	toItems := args[2:]

	// wd, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// parentDir := filepath.Dir(wd)

	// fmt.Println(parentDir)
	goPath := os.Getenv("GOPATH")
	parentDir := fmt.Sprintf("%s/src/github.com/basebandit/gocash", goPath)

	homeDir, err := homedir.Dir()

	configDir := joinPath(homeDir, ".gocash")

	if err != nil {
		boldRed.Println(err.Error())
		os.Exit(1)
	}

	//Check if config Dir exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		//Lets create the directory if it does not exist
		err := os.MkdirAll(configDir, os.FileMode(0700))
		if err != nil {
			boldRed.Println(err.Error())
			os.Exit(1)
		}
	}

	//Then we copy our currency  and config json files
	currDest := fmt.Sprintf("%s/currencies.json", configDir)
	configDest := fmt.Sprintf("%s/config.json", configDir)

	currSrc := fmt.Sprintf("%s/currencies.json", parentDir)
	configSrc := fmt.Sprintf("%s/config.json", parentDir)

	fmt.Println(currDest, currSrc)
	//copy config
	if err := copyFile(configSrc, configDest); err != nil {
		boldRed.Println(err.Error())
	}

	//copy currencies
	if err := copyFile(currSrc, currDest); err != nil {
		boldRed.Println(err.Error())
	}

	configFile := fmt.Sprintf("%s/config.json", configDir)

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

	body, err := currency.Fetch(url)
	if err != nil {
		boldRed.Println(err.Error())
		os.Exit(1)
	}

	money := new(currency.Currency)

	if rates := currency.ParseRates(body); rates != nil {
		money.Rates = rates
	}

	if base := currency.ParseBase(body); base != "" {
		money.Base = base
	}

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	//start loader
	w := wow.New(os.Stdout, spin.Get(spin.Dots), green(" Converting ..."))
	w.Start()
	time.Sleep(2 * time.Second)

	currencyFile := fmt.Sprintf("%s/currencies.json", configDir)
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
			os.Exit(1)
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

// copyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, copy the file contents from src to dst.
func copyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}

	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

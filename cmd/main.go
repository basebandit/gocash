package main

import (
	"fmt"
	"log"

	"github.com/basebandit/gocash/pkg/api"
	"github.com/basebandit/gocash/pkg/cash"
)

func main() {

	body, err := api.Fetch("http://data.fixer.io/api/latest?access_key=ea38ed8038cd5c8a41a53fa8fad29a53")
	if err != nil {
		log.Fatal(err)
	}
	if rates := cash.ParseRates(body); rates != nil {
		fmt.Println(rates)
	}
}

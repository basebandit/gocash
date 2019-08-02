package currency

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//Currency holds our base currency and the exchange rates of other currencies against the base currency
type Currency struct {
	Rates map[string]interface{}
	Base  string
}

//getExchangeRate calculates the exchange rate
func (c *Currency) getExchangeRate(to, from string) (rate float64, fxErr error) {
	// Return an error if to rate isn't in the rates array
	if _, ok := c.Rates[to]; !ok {
		fxErr = errors.New("rates: missing to rate")
		return
	}

	// Return an error if from rate isn't in the rates array
	if _, ok := c.Rates[from]; !ok {
		fxErr = errors.New("rates: missing from rate")
		return
	}

	// If `from` currency == base, return the basic exchange rate for the `to` currency
	if from == c.Base {
		to := c.Rates[to]
		rate = to.(float64)
		return
	}

	// If `to` currency == base, return the basic inverse rate of the `from` currency
	if to == c.Base {
		from := c.Rates[from]
		r := from.(float64)
		rate = 1 / r
		return
	}

	// Otherwise, return the `to` rate multiplied by the inverse of the `from` rate to get the
	// relative exchange rate between the two currencies
	t := c.Rates[to]
	rTo := t.(float64)

	f := c.Rates[from]
	rFrm := f.(float64)

	rate = rTo * (1 / rFrm)
	return
}

//Convert converts a value from one currency to another
func (c *Currency) Convert(val float64, from string, to string) (rate float64, cnvtErr error) {
	if from == "" && to == "" {
		cnvtErr = errors.New("fx error: missing either from or to rate")
		return
	}
	r, err := c.getExchangeRate(to, from)

	if err != nil {
		cnvtErr = err
		return
	}
	rate = val * r
	return
}

//UnmarshalJSON decode dynamic json data to a (key,value) pair for internal use
func unmarshalJSON(data []byte) (map[string]interface{}, error) {
	var f map[string]interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

//ParseRates retrieves rates object from the decoded json response
func ParseRates(data []byte) map[string]interface{} {
	var r map[string]interface{}
	res, err := unmarshalJSON(data)
	if err != nil {
		log.Fatal(err)
	}
	if rates, ok := res["rates"]; ok {
		r = rates.(map[string]interface{})
	}
	return r
}

//ParseBase retrieves base currency value from the decoded json response
func ParseBase(data []byte) string {
	var b string
	res, err := unmarshalJSON(data)
	if err != nil {
		log.Fatal(err)
	}
	if base, ok := res["base"]; ok {
		b = base.(string)
	}
	return b
}

//Fetch makes a GET request to the given resource url
func Fetch(url string) ([]byte, error) {
	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := httpClient.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

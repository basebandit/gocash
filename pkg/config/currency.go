package config

import (
	"encoding/json"
	"io/ioutil"
)

//GetCurrencies loads the currency json file from the specified filepath
//and returns a list of (key,value) pairs of currencies
func GetCurrencies(filepath string) (map[string]interface{}, error) {
	currFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var currency map[string]interface{}
	err = json.Unmarshal(currFile, &currency)
	if err != nil {
		return nil, err
	}
	return currency, nil
}

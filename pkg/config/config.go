package config

import (
	"encoding/json"
	"io/ioutil"
)

//Config defines the api configuration fields
type Config struct {
	ApiKey string `json:"api_key"`
	Api    string `json:"api"`
}

//LoadConfig loads the configuration file from the specified filepath
func LoadConfig(filepath string) (*Config, error) {
	//Get the config file
	configFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(configFile, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

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

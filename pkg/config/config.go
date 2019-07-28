package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ApiKey string `json:"api_key"`
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

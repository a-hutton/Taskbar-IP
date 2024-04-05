package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	PreferredIp string
}

func ReadConfig() Config {
	bytes, err := os.ReadFile("config.json")
	if err != nil {
		log.Print("Cannot read config file. Doesn't exist?")
		return Config{}
	}
	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func (config *Config) Save() {
	jsonData, err := json.Marshal(config)
	if err != nil {
		panic("Cannot convert config to json?")
	}
	err = os.WriteFile("config.json", jsonData, 0666)
	if err != nil {
		panic("cannot save config file?")
	}
}

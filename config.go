package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ListenAddr       string
	DbMySqlHost      string
	DbMySqlUser      string
	DbMySqlPassword  string
	DbMySqlName      string
	MongoHost        string
	AsteriskAddr     string
	AsteriskPort     int
	AsteriskUser     string
	AsteriskPassword string
	CrmMonFile       string
}

func NewConfig() (config *Config, err error) {
	var file []byte
	file, err = ioutil.ReadFile("./conf/config.json")

	if err != nil {
		return nil, err
	}

	config = new(Config)
	if err = json.Unmarshal(file, config); err != nil {
		return nil, err
	}

	return config, nil
}

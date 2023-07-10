package main

import (
	"encoding/json"
	"io"
	"os"
)

type AppConfig struct {
	Port     string   `json:"port"`
	DbConfig DBConfig `json:"db"`
	Static   Static   `json:"static"`
}

type Static struct {
	Port string `json:"port"`
}

type DBConfig struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Schema   string `json:"schema"`
}

func (cfg *AppConfig) InitFrom(filePath string) error {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return err
	}

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	if err := jsonFile.Close(); err == nil {
		json.Unmarshal(byteValue, &cfg)
	} else {
		return err
	}

	return nil
}

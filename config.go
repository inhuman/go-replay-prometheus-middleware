package main

import (
	"github.com/inhuman/config_merger"
	"io/ioutil"
	"log"
)

type AppConfig struct {
	UrlTypes []UrlType `json:"url_types"`
}

type UrlType struct {
	Url   string `json:"url"`
	Type  string `json:"type"`
	CType string `json:"ctype"`
}

func Config() (*AppConfig, error) {

	appConfig := &AppConfig{}

	log.SetOutput(ioutil.Discard)

	configMerger := config_merger.NewMerger(appConfig)

	configMerger.AddSource(&config_merger.JsonSource{
		Path: "config.json",
	})

	if err := configMerger.Run(); err != nil {
		return nil, err
	}
	return appConfig, nil
}

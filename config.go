package main

import (
	"fmt"
	"github.com/inhuman/config_merger"
	"os"
)

type AppConfig struct {
	UrlTypes []UrlType `json:"url_types"`
}

type UrlType struct {
	Url   string `json:"url"`
	Type  string `json:"type"`
	CType string `json:"c_type"`
}

func Config() (*AppConfig, error) {

	appConfig := &AppConfig{}

	configMerger := config_merger.NewMerger(appConfig)

	configMerger.AddSource(&config_merger.JsonSource{
		Path: "config.json",
	})

	if err := configMerger.Run(); err != nil {
		return nil, err
	}

	fmt.Fprintf(os.Stderr, "config %s\n", appConfig)

	return appConfig, nil
}
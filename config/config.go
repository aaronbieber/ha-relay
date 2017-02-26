package config

import (
	"gopkg.in/gcfg.v1"
	"log"
)

type Config struct {
	Main struct {
		Key string
	}
}

var cfg Config

func init() {
	err := gcfg.ReadFileInto(&cfg, "config.gcfg")

	if err != nil {
		log.Fatal(err)
		panic("Error loading configuration.")
	}
}

func Conf() Config {
	return cfg
}

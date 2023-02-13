package config

import (
	"flag"
	"log"
	"net/url"
	"os"
)

type Config struct {
	rosettaURL string
	networkURL string
}

func Get() *Config {
	conf := &Config{}

	flag.StringVar(&conf.rosettaURL, "rosettaurl", os.Getenv("ROSETTA_URL"), "Rosetta url")
	flag.StringVar(&conf.networkURL, "networkurl", os.Getenv("NETWORK_URL"), "Network Public url")

	flag.Parse()

	return conf
}

func (c *Config) GetRosettaUrl() (*url.URL, error) {
	u, err := url.Parse(c.rosettaURL)
	if err != nil {
		log.Fatal(err)
	}
	return u, err
}

func (c *Config) GetNetworkUrl() (*url.URL, error) {
	u, err := url.Parse(c.networkURL)
	if err != nil {
		log.Fatal(err)
	}
	return u, err
}

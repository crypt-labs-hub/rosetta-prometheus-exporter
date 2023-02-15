package config

import (
	"flag"
	"github.com/coinbase/rosetta-sdk-go/types"
	"log"
	"net/url"
	"os"
)

type Config struct {
	rosettaURL     string
	networkURL     string
	primaryNetwork *types.NetworkIdentifier
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
		log.Println(err)
	}
	return u, err
}

func (c *Config) GetNetworkUrl() (*url.URL, error) {
	u, err := url.Parse(c.networkURL)
	if err != nil {
		log.Println(err)
	}
	return u, err
}

func (c *Config) GetPrimaryNetwork() *types.NetworkIdentifier {
	return c.primaryNetwork
}

func (c *Config) SetNetwork(primaryNetwork *types.NetworkIdentifier) {
	c.primaryNetwork = primaryNetwork
}

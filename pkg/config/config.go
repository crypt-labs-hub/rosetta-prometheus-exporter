package config

import (
	"flag"
	"log"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	rosettaURL string
	networkURL string
	sampleSize string
	blockchain string
	network    string
}

func Get() *Config {
	conf := &Config{}

	flag.StringVar(&conf.rosettaURL, "rosettaurl", os.Getenv("ROSETTA_URL"), "Rosetta url")
	flag.StringVar(&conf.networkURL, "networkurl", os.Getenv("NETWORK_URL"), "Network Public url")

	flag.StringVar(&conf.sampleSize, "samplesize", os.Getenv("SAMPLE_SIZE"), "Sample size for rate calc")

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

func (c *Config) GetSampleSize() int {
	s, err := strconv.Atoi(c.sampleSize)
	if err != nil {
		log.Println(err)
		return 0
	}
	return s
}

func (c *Config) GetBlockchain() string {
	return c.blockchain
}

func (c *Config) GetNewtork() string {
	return c.network
}

func (c *Config) SetBlockchain(blockchain string) {
	c.blockchain = blockchain
}

func (c *Config) SetNetwork(network string) {
	c.network = network
}

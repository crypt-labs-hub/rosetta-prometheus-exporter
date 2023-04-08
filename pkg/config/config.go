package config

import (
	"flag"
	"log"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	// rosettaURL the url to connect to the rosetta endpoints
	rosettaURL string
	// networkURL the url to connect to normal networks besides rosetta
	networkURL string
	// sampleSize the number of samples to average the rate calculations on
	sampleSize string
	// blockchain the name of the chain used for rosetta queries
	blockchain string
	// network name used for rosetta queries
	network string
}

func Get() *Config {
	conf := &Config{}
	// ROSETTA_URL with the format http://<url>:<port>
	flag.StringVar(&conf.rosettaURL, "rosettaurl", os.Getenv("ROSETTA_URL"), "Rosetta url")
	// NETWORK_URL with the format http://<url>:<port>
	flag.StringVar(&conf.networkURL, "networkurl", os.Getenv("NETWORK_URL"), "Network Public url")
	// SAMPLE_SIZE is an integer for rate calculation
	flag.StringVar(&conf.sampleSize, "samplesize", os.Getenv("SAMPLE_SIZE"), "Sample size for rate calc")

	flag.Parse()

	return conf
}

// GetRosettaUrl rosetta url getter
func (c *Config) GetRosettaUrl() (*url.URL, error) {
	u, err := url.Parse(c.rosettaURL)
	if err != nil {
		log.Panicf("Required rosetta URL could not be parsed: %s", err)
	}
	return u, err
}

// GetNetworkUrl network url getter
func (c *Config) GetNetworkUrl() (*url.URL, error) {
	u, err := url.Parse(c.networkURL)
	if err != nil {
		log.Println(err)
	}
	return u, err
}

// GetSampleSize sample size getter
func (c *Config) GetSampleSize() int {
	s, err := strconv.Atoi(c.sampleSize)
	if err != nil {
		log.Println(err)
		// default sample size to 5
		return 5
	}
	return s
}

// GetBlockchain blockchain name getter
func (c *Config) GetBlockchain() string {
	return c.blockchain
}

// GetNetwork network name getter
func (c *Config) GetNetwork() string {
	return c.network
}

// SetBlockchain blockchain name setter
func (c *Config) SetBlockchain(blockchain string) {
	c.blockchain = blockchain
}

// SetNetwork network name setter
func (c *Config) SetNetwork(network string) {
	c.network = network
}

package main

type Config struct {
	nextUrl string
	prevUrl string
}

func getConfig() func() *Config {
	var config Config
	if config.prevUrl == "" {
		config.nextUrl = LOCATION_AREA_URL
	}
	return func() *Config {
		return &config
	}
}

func (c *Config) setNextUrl(url string) {
	c.nextUrl = url
}

func (c *Config) setPrevUrl(url string) {
	c.prevUrl = url
}

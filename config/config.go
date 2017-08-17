package config

import (
	"flag"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// FIXME(shmel1k): if possible, try to use flag

var file = flag.String("c", "conf/config.yaml", "Path to config file")

type config struct {
	Database   databaseConfig   `yaml:"database"`
	HTTPServer httpServerConfig `yaml:"http_server"`
}

type httpServerConfig struct {
	Port string `yaml:"port"`
}

type databaseConfig struct {
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Address      string `yaml:"address"`
	Port         string `yaml:"port"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

var cfg *config

func init() {
	flag.Parse()
	data, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf("failed to read file %q: %s", *file, err)
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("failed to unmarshal yaml file %q: %s", *file, err)
	}
}

func Database() databaseConfig {
	return cfg.Database
}

func HTTPServer() httpServerConfig {
	return cfg.HTTPServer
}

package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment    string `yaml:"environment" env-default:"local"` // local, dev, prod
	RepositoryFile string `yaml:"repository-file" env-required:"true"`
	HTTPServer     `yaml:"http-server"`
	Logger         `yaml:"logger"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle-timeout" env-default:"60s"`
}

type Logger struct {
	Format          string `yaml:"format" nv-default:"text"`         // text, json
	Destination     string `yaml:"destination" nv-default:"console"` // console, file
	DestinationFile string `yaml:"destination-file" nv-default:""`   // filename if destination = file
	Level           string `yaml:"level" nv-default:"debug"`         // debug, info
	AddSource       bool   `yaml:"add-source" nv-default:"false"`
}

// The function terminates the application with an error if the configuration file could not be loaded.
// In this case, we use the default logger (our main logger is not yet initialized at this stage)
func ReadConfig(path string) *Config {
	// check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}

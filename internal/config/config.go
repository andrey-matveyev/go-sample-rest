package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment    string `yaml:"environment" env-default:"local"` // local, dev, prod
	RepositoryPath string `yaml:"repository_path" env-required:"true"`
	HTTPServer     `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

// The function terminates the application with an error if the configuration file could not be loaded.
// In this case, we use the default logger (our main logger is not yet initialized at this stage)
func ReadConfig(path string) *Config {
	/*configPath, ok := os.LookupEnv(path)
	if !ok {
		log.Fatal("CONFIG_PATH is not set")
	}*/

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

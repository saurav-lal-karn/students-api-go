package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

// env-default:"production"
type Config struct {
	Env         string               `yaml:"env" env:"ENV" env-required:"true"` // Denotes env value in yaml, required = True
	StoragePath string               `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"` // Embedding the struct
}

// Must denotes the function must load and no error should be here
// Should not return err if its must load
// So that the fatal error is raised if occurs
func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		// If the config path is not inside os env, check if the param is sent via arguments in config-path
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			// Raise fatal error in case the config path is empty
			log.Fatal("Config path is not set")
		}
	}

	// Get stat error and check if the file no exists error is same
	// Which means the config path is not there in os
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file doesn't exists: %s", configPath)
	}

	var cfg Config
	// Read and serialize the config from yaml
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Cannot read config file: %s", err.Error())
	}

	return &cfg

}

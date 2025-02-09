package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type HTTPServer struct {
	Address     string        `yaml:"address"  env-default:"localhost:8080"`
	TimeOut     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeOut time.Duration `yaml:"idle_timeout" env-default:"60s"`
	//User        string        `yaml:"user" env-required:"true"`
	//Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type Storage struct {
	Host     string `yaml:"host" env-default:"localhost"`
	UserDb   string `yaml:"userdb"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port" env-default:"5432"`
	Dbname   string `yaml:"dbname"`
	SSLmode  string `yaml:"sslmode" env-default:"require"`
}

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"httpServer"`
	Storage    `yaml:"storage"`
}

func MustLoad() *Config {
	err := godotenv.Load("config/config.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("file does not exist: %s", configPath)
	}

	var cfg Config
	if err = cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	cfg.UserDb = os.Getenv("POSTGRES_USER")
	cfg.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Dbname = os.Getenv("POSTGRES_DB")

	return &cfg
}

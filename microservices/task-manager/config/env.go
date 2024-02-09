package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type Env struct {
	ServerPort       string `mapstructure:"SERVER_PORT"`
	DBUsername       string `mapstructure:"DB_USER"`
	DBPassword       string `mapstructure:"DB_PASS"`
	DBHost           string `mapstructure:"DB_HOST"`
	DBPort           string `mapstructure:"DB_PORT"`
	DBName           string `mapstructure:"DB_NAME"`
	MinioPort        string `mapstructure:"MINIO_PORT"`
	MinioUsername    string `mapstructure:"MINIO_USER"`
	MinioPassword    string `mapstructure:"MINIO_PASS"`
	MinioHost        string `mapstructure:"MINIO_HOST"`
	RabbitMQPort     string `mapstructure:"RABBITMQ_PORT"`
	RabbitMQHost     string `mapstructure:"RABBITMQ_HOST"`
	RabbitMQPassword string `mapstructure:"RABBITMQ_PASS"`
	RabbitMQUsername string `mapstructure:"RABBITMQ_USER"`
}

// NewEnv creates a new environment
func NewEnv() *Env {
	var envFile string

	if _, isPresent := os.LookupEnv("CURRENT_ENV"); isPresent {
		envFile = ".docker.env"
	} else {
		envFile = ".env"
	}

	env := Env{}
	viper.SetConfigFile(envFile)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("☠️ cannot read configuration")
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("☠️ environment can't be loaded: ", err)
	}

	return &env
}

package config

import (
	"github.com/spf13/viper"
	"log"
)

type Env struct {
	ServerPort    string `mapstructure:"SERVER_PORT"`
	DBUsername    string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASS"`
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBName        string `mapstructure:"DB_NAME"`
	MinioPort     string `mapstructure:"MINIO_PORT"`
	MinioUsername string `mapstructure:"MINIO_USER"`
	MinioPassword string `mapstructure:"MINIO_PASS"`
	MinioHost     string `mapstructure:"MINIO_HOST"`
}

// NewEnv creates a new environment
func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

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

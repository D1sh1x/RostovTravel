package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env                string     `yaml:"env"`
	JWTSecret          []byte     `yaml:"jwt_secret"`
	DBConnectionString string     `yaml:"db_connection_string"`
	HTTPServer         HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Port         string        `yaml:"port"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

var configPath string = "./config/config.yaml" //for local use

func NewConfig() *Config {

	if os.Getenv("CONFIG_PATH") != "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetTypeByDefaultValue(true)

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	config := &Config{
		Env:                v.GetString("env"),
		JWTSecret:          []byte(v.GetString("jwt_secret")),
		DBConnectionString: v.GetString("db_connection_string"),
		HTTPServer: HTTPServer{
			Port:         v.GetString("http_server.port"),
			WriteTimeout: v.GetDuration("http_server.write_timeout"),
			ReadTimeout:  v.GetDuration("http_server.read_timeout"),
			IdleTimeout:  v.GetDuration("http_server.idle_timeout"),
		},
	}

	return config
}

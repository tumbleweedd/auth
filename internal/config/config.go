package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env               string            `yaml:"env" env-default:"local"`
	EventSenderConfig EventSenderConfig `yaml:"event_sender"`

	PgConfig         PgConfig         `yaml:"postgres"`
	JWTManagerConfig JWTManagerConfig `yaml:"jwt_manager"`
	GRPCConfig       GRPCConfig       `yaml:"grpc"`
	CacheConfig      CacheConfig      `yaml:"cache"`
	KafkaConfig      KafkaConfig      `yaml:"kafka"`
	HTTPConfig       HTTPConfig       `yaml:"http"`
}

type RedisConfig struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	MinIdleConns int    `yaml:"min_idle_conns"`
	PoolSize     int    `yaml:"pool_size"`
	PoolTimeout  int    `yaml:"pool_timeout"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type JWTManagerConfig struct {
	AccessSecret  string `yaml:"accessSecret"`
	RefreshSecret string `yaml:"refreshSecret"`

	AccessTimeout  int `yaml:"accessTimeout"`  // in hours
	RefreshTimeout int `yaml:"refreshTimeout"` // in hours
}

type PgConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	SSLMode  string `yaml:"sslmode"`
}

type CacheConfig struct {
	DefaultTTL time.Duration `yaml:"default_ttl"` // in minutes
}

type EventSenderConfig struct {
	HandlePeriodMin time.Duration `yaml:"handle_period_min"`
}

type HTTPConfig struct {
	Port int `yaml:"port"`
}

type KafkaConfig struct {
	BrokerList     []string `yaml:"broker_list"`
	UserEventTopic string   `yaml:"user_event_topic"`
	Port           string   `yaml:"port"`
}

func InitConfig() *Config {
	configPath := getConfigPath()

	if configPath == "" {
		panic("config path is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("empty config path: " + err.Error())
	}

	return &cfg
}

func getConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}

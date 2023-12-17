package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Workdir     string
	Environment string `yaml:"environment" env-default:"local"`
	IsDebug     *bool  `yaml:"is_debug" env-default:"true"`
	LogLevel    string `yaml:"log_level" env-default:"debug"`
	MongoDb     `yaml:"mongodb"`
	Certificate `yaml:"certificates"`
	GameSearch  `yaml:"game_search"`
	IGBD        `yaml:"igbd"`
}
type Listen struct {
	Type   string `yaml:"type" env-default:"port"`
	BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"`
	Port   string `yaml:"port" env-default:"8080"`
}

type MongoDb struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	Database    string `yaml:"database"`
	AuthDB      string `yaml:"auth_db"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Collections string `yaml:"collections"`
}
type GameSearch struct {
	FileExtenstions []string `yaml:"extenstions"`
	Directories     []string `yaml:"directories"`
	NamesToCompare  []string `yaml:"names_to_compare"`
}

type IGBD struct {
	ID    string `yaml:"id"`
	Token string `yaml:"token"`
}

type Certificate struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

var instance *Config

var once sync.Once

func loadConfig() (string, string) {
	//search for main dir of project
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("failed to get caller information")
	}
	workdir := filepath.Dir(filename)

	for {
		if workdir == "/" {
			log.Fatal("failed to find project root")
		}
		if _, err := os.Stat(workdir + "/go.mod"); os.IsNotExist(err) {
			workdir = filepath.Dir(workdir)
		} else {
			log.Printf("workdir: %s\n", workdir)
			break
		}
	}

	configPath := workdir + "/config/local.yaml"
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file is not set: %s", configPath)
	}
	return configPath, workdir
}

func GetConfig() *Config {
	once.Do(func() {
		configPath, workdir := loadConfig()
		log.Println("read application configuration")
		instance = &Config{
			Workdir: workdir,
		}
		if err := cleanenv.ReadConfig(configPath, instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Println(help)
			log.Println(err)
		}
	})
	return instance
}

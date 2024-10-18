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
	Workdir                   string
	Environment               string `yaml:"environment" env-default:"local"`
	IsDebug                   bool   `yaml:"is_debug" env-default:"true"`
	LogLevel                  string `yaml:"log_level" env-default:"debug"`
	MongoDb                   `yaml:"mongodb"`
	Postgres                  `yaml:"postgres"`
	JWT                       `yaml:"jwt"`
	Server                    `yaml:"server"`
	Cors                      `yaml:"cors"`
	Certificates              `yaml:"certificates"`
	GameSearch                `yaml:"game_search"`
	IGBD                      `yaml:"igbd"`
	Streamer                  `yaml:"streamer"`
	UsersFileStorage          `yaml:"users_file_storage"`
	VirtualDisplayInitializer `yaml:"virtual_display_initializer"`
	UDPReader                 `yaml:"udp_reader"`
	Docker                    `yaml:"docker"`
	VideoCapture              `yaml:"video_capture"`
	AudioCapture              `yaml:"audio_capture"`
}

type Server struct {
	Type   string `yaml:"type" env-default:"port"`
	BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"`
	Port   string `yaml:"port" env-default:"8080"`
}

type Cors struct {
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	ExposedHeaders   []string `yaml:"exposed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

type Certificates struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
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

type Postgres struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Attampts int    `yaml:"attampts"`
}

type JWT struct {
	SecretKey string `yaml:"secret_key"`
	Expire    int    `yaml:"expire"`
}

type GameSearch struct {
	SystemDirectories []string `yaml:"system_directories"`
	FileExtenstions   []string `yaml:"extenstions"`
	Directories       []string `yaml:"directories"`
	NamesToCompare    []string `yaml:"names_to_compare"`
}

type IGBD struct {
	ID    string `yaml:"id"`
	Token string `yaml:"token"`
}
type Streamer struct {
	VideoCodec string `yaml:"video_codec"`
	AudioCodec string `yaml:"audio_codec"`
}

type UsersFileStorage struct {
	FsInitFilesPath string `yaml:"fs_init_files_dir"`
	Path            string `yaml:"path"`
}

type VirtualDisplayInitializer struct {
	EnableVirtualDisplaysScriptPath string `yaml:"enable_virtual_displays_script_path"`
	DisplayInfoJsonPath             string `yaml:"display_info_json_path"`
}
type UDPReader struct {
	MinPort    int `yaml:"min_port"`
	MaxPort    int `yaml:"max_port"`
	ReadBuffer int `yaml:"read_buffer"`
	UdpBuffer  int `yaml:"udp_buffer"`
}

type Docker struct {
	PulseImage   string `yaml:"pulse_image"`
	VideoImage   string `yaml:"video_image"`
	AudioImage   string `yaml:"audio_image"`
	ProtoneImage string `yaml:"protone_image"`

	CardPath       string `yaml:"card_path"`
	NetworkMode    string `yaml:"network_mode"`
	RendererPath   string `yaml:"renderer_path"`
	XauthorityPath string `yaml:"xauthority_path"`
}

type VideoCapture struct {
	Env []string `yaml:"env"`
}
type AudioCapture struct {
	Env []string `yaml:"env"`
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

package common

import (
	"genMaterials/log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

const ymlPrefix = "config"

var instance *Yaml
var mu = &sync.Mutex{}

type Yaml struct {
	Mode              string
	Echo              echo
	Database          Database
	MessageService    HmacService
	Jwt               Jwt
	RfqEmailuIRoute   HmacService
	SendEmail         bool
	BsSupplierEmailId string
}

type echo struct {
	Port   string
	Static string
}

type fluent struct {
	Path   string
	Rotate int
}

type Database struct {
	Type     string
	Host     string
	User     string
	Password string
	Name     string
	Port     string
	SslMode  string
}

type HmacService struct {
	Host             string
	Port             string
	User             string
	Scheme           string
	Secret           string
	UriPrefix        string
	AddUserUriPrefix string
}

type Jwt struct {
	Secret             string
	ExpirationTimeDays time.Duration
}

func GetYamlConfig() *Yaml {
	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		instance = &Yaml{}
	}
	return instance
}

func InitYamlConfig() *Yaml {
	yml := GetYamlConfig()
	appLogger := log.GetApplogger()
	viper.SetConfigName(ymlPrefix) // actually 'fileName + yml'
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		appLogger.Fatal("Error reading config file : ")
	}

	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			viper.Set(k, getEnvOrDefault(strings.Split(strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}"), ":")))
		}
	}

	err := viper.Unmarshal(&yml)
	if err != nil {
		appLogger.Fatal("unable to decode into struct : ")
	}
	return yml
}

func getEnvOrDefault(envList []string) string {
	env := envList[0]
	value := os.Getenv(env)
	if len(value) == 0 && len(envList) > 1 {
		envList = envList[1:]
		value = strings.Join(envList, ":")
		appLog.Debug("Using default value for key :: " + env + " value :: " + value)
	} else if len(value) == 0 {
		value = ""
		appLog.Debug("No default value found for key :: " + env + " using blank :: " + value)
	}
	return value
}

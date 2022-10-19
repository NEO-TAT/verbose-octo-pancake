package env

import (
	"github.com/spf13/viper"
	"strings"
)

var (
	IsDebugMode        bool
	RedisHosts         []string
	RedisPassword      string
	FirebaseConfigPath string
)

func init() {
	viper.AutomaticEnv()
	IsDebugMode = viper.GetBool("DEBUG")
	RedisHosts = strings.Split(viper.GetString("REDIS_HOSTS"), ",")
	RedisPassword = viper.GetString("REDIS_PASSWORD")
	FirebaseConfigPath = viper.GetString("FIREBASE_CONFIG_PATH")
}

package config

import (
	"os"
	"path/filepath"

	"github.com/croissong/releasechecker/pkg/log"
	"github.com/spf13/viper"
)

var CfgFile string
var Config configuration

func InitConfig() {
	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
		viper.SetConfigType("yaml")
	} else {
		// Find home directory.
		configDir, err := os.UserConfigDir()
		if err != nil {
			log.Logger.Error(err)
			os.Exit(1)
		}

		viper.AddConfigPath(filepath.Join(configDir, "releasechecker"))
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.ReadInConfig(); err != nil {
		log.Logger.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&Config)
	if err != nil {
		log.Logger.Fatalf("unable to decode into struct, %v", err)
	}
}

type configuration struct {
	Debug           bool
	InitDownstreams bool
	Entries         []entry
}

type entry struct {
	Name       string
	Upstream   map[string]interface{}
	Downstream map[string]interface{}
	Hooks      []map[string]interface{}
}

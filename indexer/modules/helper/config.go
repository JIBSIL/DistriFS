package helper

import (
	"fmt"

	"distrifs.dev/indexer/modules/globals"
	"github.com/spf13/viper"
)

func SetupConfig() *viper.Viper {
	fmt.Println("Setting up config file for indexer")

	supportedIndexers := make([]string, 1)
	supportedIndexers[0] = "http://example.com/server"

	viper.SetDefault("Port", 8081)

	viper.SetDefault("Debug", false)
	viper.SetDefault("HoursBetweenSyncs", 5)
	viper.SetDefault("DefaultServers", supportedIndexers)

	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			viper.SafeWriteConfigAs("./config.yaml")
		} else {
			// Config file was found but another error was produced
			panic(err)
		}
	}

	globals.Port = viper.GetInt("Port")
	globals.Debug = viper.GetBool("Debug")
	globals.TimeBetweenSyncs = viper.GetInt("HoursBetweenSyncs")
	globals.DefaultServers = viper.GetStringSlice("DefaultServers")

	return viper.GetViper()
}

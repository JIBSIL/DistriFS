package helper

import (
	"fmt"

	"distrifs.dev/server/modules/globals"
	"github.com/spf13/viper"
)

func SetupConfig() *viper.Viper {
	fmt.Println("Setting up config file for server")

	defaultAuthorizedKeys := make([]string, 1)
	defaultAuthorizedKeys[0] = "yourAuthorizedKeyHere"

	viper.SetDefault("Port", 8000)

	viper.SetDefault("OneTimeDownloads", true)
	viper.SetDefault("Directory", "./files")
	viper.SetDefault("AllFilesPassportLocked", false)
	viper.SetDefault("PassportEnabled", true)
	viper.SetDefault("AdminKeys", defaultAuthorizedKeys)
	viper.SetDefault("PassportProtectedFiles", map[string]([]string){"./files/passport/admin.txt": defaultAuthorizedKeys})

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

	globals.Options_directory = viper.GetString("Directory")
	globals.Options_oneTimeFiles = viper.GetBool("OneTimeDownloads")
	globals.Options_passport = viper.GetBool("PassportEnabled")

	return viper.GetViper()
}

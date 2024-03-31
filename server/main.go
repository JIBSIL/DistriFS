package main

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/colorstring"

	"distrifs.dev/server/modules/globals"
	"distrifs.dev/server/modules/helper"

	FileRoutes "distrifs.dev/server/routes/files"
	PassportRoutes "distrifs.dev/server/routes/passport"
)

var Badger_DB *badger.DB
var PermDB *badger.DB
var PassportDB *badger.DB

func main() {
	fmt.Printf(colorstring.Color(`
[blue][bold]DistriFS Server v%s[reset]

[blue]Distributed, internet-scale filesystem
Distributed File Server
Copyright © %d JIBSIL & Contributors

[green]Source:  https://github.com/JIBSIL/distrifs
License:  MIT%s`), globals.Version, time.Now().Year(), "\n\n")

	// set up configuration
	globals.Viper = helper.SetupConfig()

	// start webserver
	router := gin.Default()

	// open database
	db, err := badger.Open(badger.DefaultOptions("./databases/file"))

	if globals.Options_passport {
		passportdb, err := badger.Open(badger.DefaultOptions("./databases/passportdb"))

		if err != nil {
			// unrecoverable
			panic(err)
		}

		globals.PassportDB = passportdb
	}

	if err != nil {
		// unrecoverable
		panic(err)
	}

	permdb, err := badger.Open(badger.DefaultOptions("./databases/permdb"))

	if err != nil {
		// unrecoverable
		panic(err)
	}

	// setup chroot jail
	// 	err = syscall.Chroot(Options_directory)
	// 	if err != nil {
	// fmt.Println(colorstring.Color(`
	// 		[red][bold][underline]SECURITY WARNING
	// [reset][red]An environment that does not support sandboxing was detected.
	// On Linux: Run this program with sudo
	// On Windows: You can ignore this - but please consider running on a Linux virtual machine, as it supports more advanced security features

	// Following the steps above is recommended but not required. Without sandboxing, your system is more vulnerable.

	// Path sanitization will be enabled for some security.
	// 		`))
	// 		InternalOptions_pathSanitization = true
	// 	}

	globals.Badger_DB = db
	globals.PermDB = permdb

	// identify server
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"alive":   true,
			"server":  "DistriFS Distributed File Server",
			"version": globals.Version,
		})
	})

	// setup storage endpoints
	files := router.Group("/files")
	{
		files.GET("/dir", FileRoutes.ReadDirRoute)
		files.GET("/read", FileRoutes.ReadFileRoute)
		files.GET("/stat", FileRoutes.StatFileRoute)
		files.GET("/download", FileRoutes.DownloadFileRoute)
		files.GET("/exists", FileRoutes.ExistsRoute)
	}

	passport := router.Group("/passport")
	{
		passport.POST("/verify", PassportRoutes.VerifyPassportRoute)
		passport.GET("/authenticate", PassportRoutes.AuthenticateRoute)
		passport.POST("/getKey", PassportRoutes.GetKeyRoute)
	}

	address := "0.0.0.0:" + fmt.Sprint(globals.Viper.GetInt("Port"))

	router.Run(address) // listen and serve on 0.0.0.0:port
}

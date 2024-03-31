package main

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/colorstring"

	"distrifs.dev/indexer/modules/globals"
	"distrifs.dev/indexer/modules/helper"

	Routes "distrifs.dev/indexer/routes"
)

func main() {
	// first, print version info
	fmt.Printf(colorstring.Color(`
[blue][bold]DistriFS Indexer v%s[reset]

[blue]Distributed, internet-scale filesystem
Distributed Indexer Server
Copyright © %d JIBSIL & Contributors

[green]Source:  https://github.com/JIBSIL/distrifs
License:  MIT%s`), globals.Version, time.Now().Year(), "\n\n")
	// start webserver
	router := gin.Default()

	globals.Viper = helper.SetupConfig()

	if globals.Debug {
		helper.CheckDummyServer()
	}

	// open database
	indexer, err := badger.Open(badger.DefaultOptions("./indexdb"))

	if err != nil {
		// unrecoverable
		panic(err)
	}

	globals.Badger_DB = indexer
	globals.IndexerData = make(globals.Indexer)

	// add default indexers
	for _, indexer := range globals.DefaultServers {
		_, ok := globals.IndexerData[indexer]

		if !ok {
			server := globals.IndexerServer{}
			server.Online = true
			server.Hashes = 0
			server.Files = globals.HashItem{}

			globals.IndexerData[indexer] = server
		}
	}

	helper.RunFullSync()

	// setup storage endpoints
	files := router.Group("/indexer")
	{
		files.GET("/get", Routes.GetFileRoute)
		files.GET("/search", Routes.SearchFileRoute)
		files.GET("/suggest", Routes.PutFileRoute)
		files.GET("/getAll", Routes.GetAllHelpRoute)
		files.GET("/getAll/indexers", Routes.GetAllFormatted)
		files.GET("/getAll/files", Routes.GetAllFilesRoute)
	}

	// dev option
	if globals.Debug {
		helper.GenerateDummyTree()
	}

	router.Run(":" + fmt.Sprint(globals.Port)) // listen and serve on 0.0.0.0:8081
}

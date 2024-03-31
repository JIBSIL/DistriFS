package globals

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/spf13/viper"
)

// this can stay how it is (hardcoded)
var Version = "0.0.1-DEV"

var Badger_DB *badger.DB
var Viper *viper.Viper

var IndexerData Indexer

var TimeBetweenSyncs = 5 // 5 hours
var Port = 8081
var Debug = true
var DefaultServers []string

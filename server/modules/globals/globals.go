package globals

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/spf13/viper"
)

// this can stay how it is (hardcoded)
// but everything else has been moved into config.yaml
var Version = "0.0.1-DEV"

var Options_oneTimeFiles bool
var Options_directory string
var Options_passport bool

var Badger_DB *badger.DB
var PermDB *badger.DB
var PassportDB *badger.DB
var Viper *viper.Viper

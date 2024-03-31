package globals

import (
	"time"
)

type HashItem struct {
	FileName    string
	Size        int64
	IsDirectory bool
	ModTime     time.Time
	Hash        string
	SubFiles    map[string]HashItem
}

type IndexerServer struct {
	Online   bool
	Hashes   int
	Files    HashItem
	LastSync time.Time
}

type Indexer map[string]IndexerServer

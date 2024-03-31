package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"distrifs.dev/indexer/modules/globals"
)

// this package generates a dummy tree so sorting can be tested

type recursiveSortResponse struct {
	error  bool
	data   globals.HashItem
	hashes int
}

// make error handling a bit easier
func recursiveSortError(error bool) recursiveSortResponse {
	response := recursiveSortResponse{}
	response.error = error
	return response
}

func recursiveSort(pathInput string, server globals.IndexerServer) recursiveSortResponse {
	var hashes int
	err := filepath.WalkDir(pathInput,
		func(path string, info os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			fileinfo, _ := info.Info()

			pathFormatted, _ := filepath.Split(path)
			directories := strings.Split(path, "/")

			fmt.Println(info.Name(), fileinfo.Size())
			fmt.Println(directories, pathFormatted)

			dirRelative := server.Files

			for _, value := range directories {
				val, ok := dirRelative.SubFiles[value]
				if !ok {
					// current doesn't exist
					fmt.Println("adding directory " + value)

					if info.IsDir() {
						dirRelative.SubFiles[value] = globals.HashItem{
							FileName:    info.Name(),
							IsDirectory: true,
							SubFiles:    make(map[string]globals.HashItem),
						}
					} else {
						checksum, err := ChecksumFile(path)

						if err != nil {
							panic(err)
						}
						dirRelative.SubFiles[value] = globals.HashItem{
							FileName:    info.Name(),
							Size:        fileinfo.Size(),
							IsDirectory: false,
							Hash:        checksum,
							ModTime:     fileinfo.ModTime(),
							SubFiles:    make(map[string]globals.HashItem),
						}

						hashes++
					}
				} else {
					dirRelative = val
				}
			}

			//server.Files[]

			return nil
		})

	if err != nil {
		return recursiveSortError(false)
	}

	return recursiveSortResponse{
		error:  false,
		data:   server.Files,
		hashes: hashes,
	}
}

func GenerateDummyTree() {
	fmt.Println("Running dummy tree generation (to test indexer)")
	// create a new indexer server
	server := globals.IndexerServer{}

	server.Online = true
	server.Hashes = 0
	server.Files = globals.HashItem{}
	server.Files.SubFiles = make(map[string]globals.HashItem)

	res := recursiveSort("./example_file_structure", server)

	if !res.error {
		indexer := globals.IndexerData["dummy"]

		indexer.Online = true
		indexer.Hashes = res.hashes
		indexer.Files = res.data

		fmt.Println("SUCCESS! Dummy tree successfully generated!")

		globals.IndexerData["dummy.local"] = indexer
	} else {
		fmt.Println("ERROR! Dummy tree was not generated!")
	}
}

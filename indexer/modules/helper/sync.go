package helper

// file hosting all of the sync functions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"distrifs.dev/indexer/modules/globals"
)

type FileInfo struct {
	Name        string
	IsDirectory bool
}

type FileResponse struct {
	Success bool       `json:"success"`
	Data    []FileInfo `json:"data"`
}

type ServerResponse struct {
	Error bool
	Data  globals.IndexerServer
}

type LoopDirectoryReturn struct {
	Data   map[string]globals.HashItem
	Hashes int
}

type StatFileReturn struct {
	Success      bool      `json:"success"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	IsDirectory  bool      `json:"isDirectory"`
	LastModified time.Time `json:"lastModified"`
	Hash         string    `json:"hash"`
}

func CheckDummyServer() {
	fmt.Println("Checking localhost server...")
	var dummy = "http://localhost:8000"

	CheckServer(dummy, globals.IndexerServer{})
}

func RunFullSync() {
	fmt.Println("Running full sync on all servers")

	for server_uri := range globals.IndexerData {
		server := globals.IndexerData[server_uri]

		if !server.LastSync.Add(time.Duration(globals.TimeBetweenSyncs) * time.Hour).Before(time.Now()) {
			// not enough time has passed!
			return
		}

		fmt.Println("Indexer is out of sync! Checking:", server_uri)

		response := CheckServer(server_uri, server)

		if !response.Error {
			response.Data.LastSync = time.Now()
			globals.IndexerData[server_uri] = response.Data
		} else {
			fmt.Println("Error syncing server " + server_uri)
		}
	}
}

func AppendTree(uri string) globals.IndexerServer {
	// assumes that there is no existing server

	fmt.Println("Checking new server " + uri)

	server := globals.IndexerServer{}

	server.Online = true
	server.Hashes = 0
	server.Files = globals.HashItem{}
	server.Files.SubFiles = make(map[string]globals.HashItem)
	server.LastSync = time.Now()

	response := CheckServer(uri, server)

	if response.Error {
		fmt.Println("Server is dead!")
		server.Online = false
	}

	server = response.Data
	return server
}

func GetHashInfoFromServer(server string, path string) StatFileReturn {
	var res, err = http.Get((server + "/files/stat?filename=" + path))

	if err != nil {
		fmt.Println("Server died while querying hash! Returning")
		return StatFileReturn{}
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return StatFileReturn{}
	}

	var statReturn StatFileReturn

	json.Unmarshal(body, &statReturn)

	return statReturn
}

func LoopDirectory(server string, dir string, subFiles map[string]globals.HashItem) LoopDirectoryReturn {

	// get request with query params
	if dir == "" {
		dir = "/"
	}

	if subFiles == nil {
		subFiles = make(map[string]globals.HashItem)
	}

	var res, err = http.Get((server + "/files/dir?dirname=" + dir))

	var hashes = 0
	var returned = LoopDirectoryReturn{}

	if err != nil {
		fmt.Println("Server died while looping! Returning..")
		return returned
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println("Error reading body, returning..")
		return returned
	}

	var parsed FileResponse
	json.Unmarshal(body, &parsed)

	if dir == "/" {
		dir = ""
	}

	if parsed.Success {
		for _, element := range parsed.Data {
			var new_directory = dir + "/" + element.Name
			if element.IsDirectory {
				fmt.Println("Directory: ", new_directory)

				response := LoopDirectory(server, new_directory, make(map[string]globals.HashItem))
				hashes += response.Hashes

				entry := globals.HashItem{
					FileName:    element.Name,
					IsDirectory: true,
					SubFiles:    response.Data,
				}

				hashinfo := GetHashInfoFromServer(server, new_directory)

				entry.ModTime = hashinfo.LastModified

				subFiles[element.Name] = entry

			} else {
				fmt.Println("File: ", element.Name)

				entry := globals.HashItem{
					FileName: element.Name,
					// Size:        fileinfo.Size(),
					IsDirectory: false,
					// Hash:        checksum,
					// ModTime:     fileinfo.ModTime(),
					SubFiles: make(map[string]globals.HashItem),
				}

				hashinfo := GetHashInfoFromServer(server, new_directory)

				entry.Size = hashinfo.Size
				entry.Hash = hashinfo.Hash
				entry.ModTime = hashinfo.LastModified

				subFiles[element.Name] = entry
				hashes++
			}
		}
	}

	returned.Hashes = hashes
	returned.Data = subFiles

	return returned
}

func CheckServer(server string, indexer_server globals.IndexerServer) ServerResponse {
	// ping server http
	// first ping is jsut to check the server is up and reading main dir works

	var res, err = http.Get((server + "/files/dir?dirname=/"))

	var response = ServerResponse{}
	response.Error = true

	if err != nil {
		fmt.Println(err)
		fmt.Println("Server is not alive, skipping")
		return response
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println("Error reading body, returning..")
		return response
	}

	var parsed FileResponse
	var loopFiles map[string]globals.HashItem

	json.Unmarshal(body, &parsed)

	startTime := time.Now()

	if parsed.Success {
		directoryResponse := LoopDirectory(server, "", indexer_server.Files.SubFiles)

		loopFiles = directoryResponse.Data
		indexer_server.Hashes = directoryResponse.Hashes
	} else {
		fmt.Println("Server is not alive, skipping")
	}

	duration := time.Since(startTime).Milliseconds()
	fmt.Println("Full server index done in", duration, "ms")

	indexer_server.Files.SubFiles = loopFiles

	response.Error = false
	response.Data = indexer_server

	return response
}

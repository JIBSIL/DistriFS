package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

type GenericResponse struct {
	Data    string `json:"data"`
	Success bool   `json:"success"`
}

type DumpIndexersResponse struct {
	Data    []string `json:"data"`
	Success bool     `json:"success"`
}

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

type GetAllFilesResponse struct {
	Data    map[string]IndexerServer `json:"data"`
	Success bool                     `json:"success"`
}

type SearchFileResponse struct {
	Data    []FileTarget `json:"data"`
	Success bool         `json:"success"`
}

type GetFileResponse struct {
	Data struct {
		FileName    string `json:"FileName"`
		Size        int    `json:"Size"`
		IsDirectory bool   `json:"IsDirectory"`
		ModTime     string `json:"ModTime"`
		Hash        string `json:"Hash"`
		SubFiles    struct {
		} `json:"SubFiles"`
	} `json:"data"`
	Success bool `json:"success"`
}

func SendGetRequest(t *testing.T, url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		t.Errorf("Failed to get send request: %s", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response: %s", err)
	}

	return body
}

func SendPostRequest(t *testing.T, url string, reqBody interface{}) []byte {
	json_data, err := json.Marshal(reqBody)
	if err != nil {
		t.Errorf("Failed to marshal POST request: %s", err)
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		t.Errorf("Failed to get send request: %s", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response: %s", err)
	}

	return body
}

// this tests both files/read and files/download
func TestDumpIndexers(t *testing.T) {
	// Run download tests
	response := SendGetRequest(t, "http://localhost:8081/indexer/getAll/indexers")
	var result DumpIndexersResponse
	if err := json.Unmarshal(response, &result); err != nil {
		t.Errorf("Failed to unmarshal JSON: %s", err)
	}

	if result.Success {
		if len(result.Data) == 1 && result.Data[0] == "http://localhost:8000" {
			t.Log("Dump test passed")
		}
	} else {
		t.Errorf("Dump test failed")
	}
}

func TestGetAllFiles(t *testing.T) {
	// Test if hello.txt exists (lightweight Exists route)
	response := SendGetRequest(t, "http://localhost:8081/indexer/getAll/files")
	var result GetAllFilesResponse
	if err := json.Unmarshal(response, &result); err != nil {
		t.Errorf("Failed to unmarshal JSON: %s", err)
	}

	if result.Success {
		if result.Data["http://localhost:8000"].Online == true && result.Data["http://localhost:8000"].Hashes == 3 {
			t.Log("GetAllFiles test passed")
		} else {
			t.Errorf("GetAllFiles test failed (2)")
		}
	} else {
		t.Errorf("GetAllFiles test failed")
	}
}

func TestGetFile(t *testing.T) {
	// Test if ReadDir in main dir is functioning
	response := SendGetRequest(t, "http://localhost:8081/indexer/get?server=http://localhost:8000&hash=1a9ef7e5588a84b048a7f60468c270888dada51adfd1f9276662d4a275e73198")
	var result GetFileResponse
	if err := json.Unmarshal(response, &result); err != nil {
		t.Errorf("Failed to unmarshal JSON: %s", err)
	}

	if result.Success {
		if result.Data.FileName == "hello.txt" && result.Data.Size == 15 {
			t.Log("GetFile test passed")
		} else {
			t.Errorf("GetFile test failed (wrong size of element, is indexer running default config?)")
		}
	} else {
		t.Errorf("GetFile test failed")
	}
}

func TestSearchFile(t *testing.T) {
	// Test if StatFile in main dir is functioning
	response := SendGetRequest(t, "http://localhost:8081/indexer/search?server=http://localhost:8000&query=hello.txt")
	var result SearchFileResponse
	if err := json.Unmarshal(response, &result); err != nil {
		t.Errorf("Failed to unmarshal JSON: %s", err)
	}

	if result.Success {
		if result.Data[0].PercentMatched == 1 && result.Data[0].Hash == "1a9ef7e5588a84b048a7f60468c270888dada51adfd1f9276662d4a275e73198" { // hash of Hello DistriFS!
			t.Log("SearchFile test passed")
		} else {
			t.Error("SearchFile test failed (wrong hash or % matched, is server running default config?)")
		}
	} else {
		t.Errorf("SearchFile test failed")
	}
}

func TestPutFile(t *testing.T) {
	response := SendGetRequest(t, "http://localhost:8081/indexer/suggest?server=http://localhost:8000")
	var result GenericResponse
	if err := json.Unmarshal(response, &result); err != nil {
		t.Errorf("Failed to unmarshal JSON: %s", err)
	}

	if result.Success || result.Data == "server_already_exists" {
		t.Log("PutFile test passed")
	} else {
		t.Errorf("PutFile test failed")
	}
}

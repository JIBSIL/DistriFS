package files

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type GenericResponse struct {
	Data    string `json:"data"`
	Success bool   `json:"success"`
}

type SimpleGenericResponse struct {
	Success bool `json:"success"`
}

type ReadDirResponse struct {
	Data []struct {
		Name        string `json:"Name"`
		IsDirectory bool   `json:"IsDirectory"`
	} `json:"data"`
	Success bool `json:"success"`
}

type StatFileResponse struct {
	Hash         string `json:"hash"`
	IsDirectory  bool   `json:"isDirectory"`
	LastModified string `json:"lastModified"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Success      bool   `json:"success"`
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
func TestDownload(t *testing.T) {
	// Run download tests
	response := SendGetRequest(t, "http://localhost:8000/files/read?filename=hello.txt")
	var result GenericResponse
	if err := json.Unmarshal(response, &result); err != nil {
		t.Errorf("Failed to unmarshal JSON: %s", err)
	}

	response = SendGetRequest(t, "http://localhost:8000/files/download?id="+result.Data)
	t.Log(string(response))

	if string(response) == "Hello DistriFS!" {
		t.Log("Download test passed")
	} else {
		t.Errorf("Download test failed")
	}
}

func TestExists(t *testing.T) {
	// Test if hello.txt exists (lightweight Exists route)
	response := SendGetRequest(t, "http://localhost:8000/files/exists?filename=hello.txt")
	var result SimpleGenericResponse
	if err := json.Unmarshal(response, &result); err != nil {
		t.Errorf("Failed to unmarshal JSON: %s", err)
	}

	if result.Success {
		t.Log("Exists test passed")
	} else {
		t.Errorf("Exists test failed")
	}
}

func TestReadDir(t *testing.T) {
	// Test if ReadDir in main dir is functioning
	response := SendGetRequest(t, "http://localhost:8000/files/dir?dirname=.")
	var result ReadDirResponse
	if err := json.Unmarshal(response, &result); err != nil {
		t.Errorf("Failed to unmarshal JSON: %s", err)
	}

	if result.Success {
		if len(result.Data) == 3 && result.Data[0].Name == "hello.txt" {
			t.Log("ReadDir test passed")
		} else {
			t.Errorf("ReadDir test failed (wrong len of elements or wrong element 0, is server running default config?)")
		}
	} else {
		t.Errorf("ReadDir test failed")
	}
}

func TestStatFile(t *testing.T) {
	// Test if StatFile in main dir is functioning
	response := SendGetRequest(t, "http://localhost:8000/files/stat?filename=hello.txt")
	var result StatFileResponse
	if err := json.Unmarshal(response, &result); err != nil {
		t.Errorf("Failed to unmarshal JSON: %s", err)
	}

	if result.Success {
		if result.Hash == "1a9ef7e5588a84b048a7f60468c270888dada51adfd1f9276662d4a275e73198" { // hash of Hello DistriFS!
			t.Log("StatFile test passed")
		} else {
			t.Errorf("StatFile test failed (wrong hash, is server running default config?)")
		}
	} else {
		t.Errorf("StatFile test failed")
	}
}

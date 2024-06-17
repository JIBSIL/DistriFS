package files

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
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

func SendGetRequestNoOutputNoTesting(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		// we silently fail in this function (used in benchmark)
		return []byte{}
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}
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

func FileExists(file string) bool {
	// Test if hello.txt exists (lightweight Exists route)
	response := SendGetRequestNoOutputNoTesting("http://localhost:8000/files/exists?filename=" + file)
	var result SimpleGenericResponse
	if err := json.Unmarshal(response, &result); err != nil {
		fmt.Printf("Failed to unmarshal JSON: %s", err)
	}

	if result.Success {
		return true
	} else {
		return false
	}
}

func RunDownloadBenchmark(b *testing.B, file string) bool {

	filename := strings.ReplaceAll(file, "./benchmark/", "")

	f, err := os.Create("./tmp/" + filename + "-downloaded-" + fmt.Sprint(rand.Int()) + ".tmp")
	if err != nil {
		b.Errorf("Failed to create tmp file: %s", err)
		return false
	}
	defer f.Close()

	response := SendGetRequestNoOutputNoTesting("http://localhost:8000/files/read?filename=" + file)
	var result GenericResponse
	if err := json.Unmarshal(response, &result); err != nil {
		b.Errorf("Failed to unmarshal JSON: %s", err)
	}

	response = SendGetRequestNoOutputNoTesting("http://localhost:8000/files/download?id=" + result.Data)

	_, err = io.Copy(f, bytes.NewReader(response))

	return err == nil
}

func ThreadpoolDownloaderWorker(id int, jobs <-chan int, results chan<- int, file string, b *testing.B) {
	for j := range jobs {
		if !RunDownloadBenchmark(b, file) {
			b.Error("Failed to download " + file)
		}
		results <- j
	}
}

func ThreadpoolDownload(b *testing.B, file string, num int) {
	jobs := make(chan int, num)
	results := make(chan int, num)

	for w := 1; w <= 5; w++ {
		go ThreadpoolDownloaderWorker(w, jobs, results, file, b)
	}

	for j := 1; j <= num; j++ {
		jobs <- j
	}
	close(jobs)

	// collect results
	for a := 1; a <= num; a++ {
		<-results
	}
}

func BenchmarkDownload(b *testing.B) {
	// test if file exists to server
	if FileExists("./benchmark/1gb-file.bin") && FileExists("./benchmark/100mb-file.bin") && FileExists("./benchmark/1mb-file.bin") {
		initial := time.Now()
		// create tmp folder
		if _, err := os.Stat("./tmp"); os.IsNotExist(err) {
			os.Mkdir("./tmp", 0755)
		}

		// download one 1gb file to test
		if !RunDownloadBenchmark(b, "./benchmark/1gb-file.bin") {
			b.Errorf("Failed to download 1gb-file.bin")
		}
		b.Logf("SEQUENTIAL: Downloaded 1gb-file.bin in %s (%s MB/s)", time.Since(initial), fmt.Sprint(1024/int(time.Since(initial).Seconds())))

		// download ten 100mb files
		initial = time.Now()
		for i := 0; i < 10; i++ {
			if !RunDownloadBenchmark(b, "./benchmark/100mb-file.bin") {
				b.Errorf("Failed to download 100mb-file.bin")
			}
		}
		b.Logf("SEQUENTIAL: Downloaded 10 100mb-file.bin files sequentially in %s (%s MB/s)", time.Since(initial), fmt.Sprint((100*10)/int(time.Since(initial).Seconds())))

		// download one hundred 1mb files
		initial = time.Now()
		for i := 0; i < 100; i++ {
			if !RunDownloadBenchmark(b, "./benchmark/1mb-file.bin") {
				b.Errorf("Failed to download 1mb-file.bin")
			}
		}
		b.Logf("SEQUENTIAL: Downloaded 100 1mb-file.bin files sequentially in %s (%s MB/s)", time.Since(initial), fmt.Sprint((1*100)/float32(time.Since(initial).Seconds())))

		// run parallel tests
		initial = time.Now()
		ThreadpoolDownload(b, "./benchmark/100mb-file.bin", 10)
		b.Logf("PARALLEL: Downloaded 10 100mb-file.bin files in %s (%s MB/s)", time.Since(initial), fmt.Sprint((100*10)/int(time.Since(initial).Seconds())))

		initial = time.Now()
		ThreadpoolDownload(b, "./benchmark/1mb-file.bin", 100)
		b.Logf("PARALLEL: Downloaded 100 1mb-file.bin files in %s (%s MB/s)", time.Since(initial), fmt.Sprint((1*100)/float32(time.Since(initial).Seconds())))

		// delete tmp folder
		os.RemoveAll("./tmp")
	} else {
		b.Errorf("Make sure to run benchmarks/generate script and move the generated files into the benchmark folder before running this benchmark! Also make sure the server is running.")
	}
}

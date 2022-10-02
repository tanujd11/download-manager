package downloader

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestDownload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(".", "testdata", "1gig.bin"))
	}))
	defer server.Close()
	downloadOptions := DownloadOptions{
		DownloadPath: "1gig.bin",
		NumConcParts: 15,
		Workers:      5,
	}
	downloader := Downloader{
		DownloadOptions: downloadOptions,
		Progress:        0,
		Mutex:           &sync.Mutex{},
	}

	downloadPath, err := downloader.Download(server.URL)
	if err != nil {
		t.Errorf("error downloading the file: %s", err.Error())
	}

	_, err = os.ReadFile(downloadPath)

	if err != nil {
		t.Errorf("file not present at downloadPath: %s", err.Error())
	}

	f1, err := os.Stat(filepath.Join(".", "testdata", "1gig.bin"))
	if err != nil {
		return
	}
	f2, err := os.Stat("1gig.bin")
	if err != nil {
		return
	}

	if f1.Size() != f2.Size() {
		t.Errorf("downloaded file not the same as server file")
	}
}

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

	t.Logf("download succesful")
	defer server.Close()
}

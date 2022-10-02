package downloader

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/tanujd11/download-manager/internal/chunk"
)

const (
	chunkDir = "/tmp"
)

// DownloadClient is a simple HTTP Downloader that supports
// concurrent downloading of large files.
type DownloadClient interface {
	// Download concurrently downloads the given url into the configured downloadDir using
	// DownloadOptions.NumConcParts
	// appropriately and returns paths to the locally downloaded files or error.
	Download(fileUrl string) (downloadPath string, err error)
	// GetDownloadProgress obtains the download progress percentage fileUrl rounded down to nearest integer
	GetDownloadProgress(fileUrl string) int
}

// Downloader is a struct which handles all the downloading information and concurrency control
type Downloader struct {
	DownloadOptions DownloadOptions
	Progress        int
	Mutex           *sync.Mutex
}

// Download function downloads the file
func (d *Downloader) Download(fileUrl string) (downloadPath string, err error) {
	// Create a head request to get the content length of the
	res, err := http.Head(fileUrl)
	if err != nil {
		return "", err
	}
	contentLength := int(res.ContentLength)
	if err != nil {
		return "", err
	}

	chunkCount := d.DownloadOptions.NumConcParts
	chunkSize := contentLength / chunkCount
	chunks := []chunk.Chunk{}

	// Create chunks according to the number of concurrent Parts
	for i := 0; i < chunkCount; i++ {
		if i < chunkCount-1 {
			chunks = append(chunks, chunk.Chunk{Start: i * chunkSize, End: (i+1)*chunkSize - 1})
		} else {
			chunks = append(chunks, chunk.Chunk{Start: i * chunkSize, End: contentLength - 1})
		}

		chunks[i].SetIndex(i)
	}

	// Download all the chunks created
	wg := sync.WaitGroup{}
	for _, c := range chunks {
		wg.Add(1)
		go func(c chunk.Chunk) {
			defer wg.Done()
			err := c.Download(fileUrl, chunkDir)
			if err != nil {
				panic(err)
			}
			d.Mutex.Lock()
			d.Progress = d.Progress + 100/chunkCount
			d.Mutex.Unlock()
		}(c)
	}

	wg.Wait()

	//merge the downloaded temp files to get the final output
	f, err := os.OpenFile(d.DownloadOptions.DownloadPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for i := range chunks {
		fileBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/file-%v.tmp", chunkDir, i))
		if err != nil {
			return "", err
		}

		if _, err = f.Write(fileBytes); err != nil {
			return "", err
		}
	}

	// removing chunks as download succesful
	err = chunk.Cleanup(chunkDir)

	return d.DownloadOptions.DownloadPath, nil
}

func (d *Downloader) GetDownloadProgress(fileUrl string) int {
	d.Mutex.Lock()
	progress := d.Progress
	d.Mutex.Unlock()
	return progress
}

type DownloadOptions struct {
	DownloadPath string
	// NumConcParts represents max number of go-routines used to download diff parts
	// of a large file simultaneously.
	NumConcParts int
}

func NewDownloadClient(opts DownloadOptions) DownloadClient {
	downloadClient := &Downloader{DownloadOptions: opts}
	return downloadClient
}

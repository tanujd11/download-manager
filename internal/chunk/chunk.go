package chunk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type Chunk struct {
	Start int
	End   int
	Index int
}

func (chunk *Chunk) SetStart(start int) *Chunk {
	chunk.Start = start
	return chunk
}

func (chunk *Chunk) SetEnd(end int) *Chunk {
	chunk.End = end
	return chunk
}

func (chunk *Chunk) SetIndex(index int) *Chunk {
	chunk.Index = index
	return chunk
}

func (chunk Chunk) Download(fileUrl string, chunkDir string) error {
	req, err := http.NewRequest("GET", fileUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", chunk.Start, chunk.End))
	res, err := http.DefaultClient.Do(req)
	if res.StatusCode > 299 {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	ioutil.WriteFile(fmt.Sprintf("%s/file-%v.tmp", chunkDir, chunk.Index), body, 0644)
	return nil
}

func Merge(chunks []Chunk, downloadPath string, chunkDir string) error {
	f, err := os.OpenFile(downloadPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for i := range chunks {
		fileBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/file-%v.tmp", chunkDir, i))
		if err != nil {
			return err
		}

		if _, err = f.Write(fileBytes); err != nil {
			return err
		}
	}

	return nil
}

func Cleanup(chunkDir string) error {
	files, err := filepath.Glob(fmt.Sprintf("%s/file-*.tmp", chunkDir))
	if err != nil {
		return err
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}

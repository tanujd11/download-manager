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

func (chunk *Chunk) Download(fileUrl string, chunkDir string) error {
	req, err := http.NewRequest("GET", fileUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", chunk.Start, chunk.End))
	res, _ := http.DefaultClient.Do(req)
	if err != nil {
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
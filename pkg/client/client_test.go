package client

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"testing"
)

func TestDownloadRange(t *testing.T) {
	url := "https://civitai.com/api/download/models/50722"
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	log.Printf("%s\n", resp.Header.Get("Content-Disposition"))
	log.Printf("%s\n", resp.Header.Get("Content-Length"))
	contentLength := resp.Header.Get("Content-Length")
	fileSize, err := strconv.Atoi(contentLength)
	if err != nil {
		log.Println(err)
		return
	}
	chunkSize := 10 * 1024 * 1024
	cnt := fileSize / chunkSize
	start := cnt * chunkSize
	end := (cnt+1)*chunkSize - 1
	if end > fileSize-1 {
		end = fileSize - 1
	}
	log.Println("===", cnt)
	log.Println("===start", start)
	log.Println("===end", end)

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	resp2, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp2.Body.Close()
	buf := make([]byte, chunkSize)
	n, err := io.CopyN(bytes.NewBuffer(buf), resp2.Body, int64(chunkSize))
	_ = buf[:n]
	log.Println("===n", n)
	log.Println("===err", err)
}

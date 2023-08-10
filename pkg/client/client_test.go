package client

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestDownloadRange(t *testing.T) {
	url := "https://civitai.com/api/download/models/50722"
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", 0, 100*1024*1024-1))
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	log.Printf("%s\n", resp.Header.Get("Content-Disposition"))
	log.Printf("%s\n", resp.Header.Get("Content-Length"))
}

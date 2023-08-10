package client

import (
	"io"
	"log"
	"testing"
)

func TestDownloadRange(t *testing.T) {
	url := "https://civitai.com/api/download/models/42985"
	var aliClient = &AliClient{chunkSize: 10 * 1024 * 1024}
	err := aliClient.getFileMeta(url)
	if err != nil {
		log.Println(err)
		return
	}

	totalIter := aliClient.fileSize / aliClient.chunkSize
	log.Println("===fileSize", aliClient.fileSize)
	log.Println("===chunkSize", aliClient.chunkSize)
	log.Println("===totalIter", totalIter)
	cnt := 0
	iter := 0
	for i := 0; i <= totalIter; i++ {
		start := i * aliClient.chunkSize
		end := (i+1)*aliClient.chunkSize - 1
		if end > aliClient.fileSize-1 {
			end = aliClient.fileSize - 1
		}
		data, err := aliClient.downloadRange(url, start, end)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
				return
			}
		}
		iter += 1
		cnt += len(data)

		log.Println("===current", iter)
		log.Println("===byte read", len(data))
		log.Println("===byte total", cnt)
	}

	log.Println("[success] result:")
	log.Printf("totalIter: %d, iter: %d\n", totalIter, iter)
	log.Printf("fileSize: %d, bytes send: %d\n", aliClient.fileSize, cnt)
}

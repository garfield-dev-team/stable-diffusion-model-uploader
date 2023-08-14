package client

import (
	"bytes"
	"context"
	"crypto/md5"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetFileSize(t *testing.T) {
	url := "http://localhost:8066/vendors-node_modules_react-dom_client_js-node_modules_webpack-dev-server_client_index_js_prot-8907fc.js"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("===len", len(data))
	log.Printf("===sum: %x\n", md5.Sum(data))
}

func TestDownloadChunk(t *testing.T) {
	url := "http://localhost:8066/vendors-node_modules_react-dom_client_js-node_modules_webpack-dev-server_client_index_js_prot-8907fc.js"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	buf := make([]byte, 1024)
	res := make([]byte, 0)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		res = append(res, buf[:n]...)
	}
	log.Printf("===sum: %x\n", md5.Sum(res))
}

func TestDownloadRange(t *testing.T) {
	_, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()
	//url := "https://civitai.com/api/download/models/111612"
	url := "http://localhost:8066/vendors-node_modules_react-dom_client_js-node_modules_webpack-dev-server_client_index_js_prot-8907fc.js"
	//var aliClient = &AliClient{chunkSize: 10 * 1024 * 1024}
	var aliClient = &AliClient{chunkSize: 1024, fileSize: 1188339}
	checksum := "0f3844283d8fdbb7a21b337b91457c1e"
	//err := aliClient.getFileMeta(url)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	totalIter := aliClient.fileSize / aliClient.chunkSize
	log.Println("===fileSize", aliClient.fileSize)
	log.Println("===chunkSize", aliClient.chunkSize)
	log.Println("===totalIter", totalIter)
	var buf bytes.Buffer
	cnt := 0
	iter := 0
	for i := 0; i <= totalIter; i++ {
		start := i * aliClient.chunkSize
		end := (i+1)*aliClient.chunkSize - 1
		if end > aliClient.fileSize-1 {
			end = aliClient.fileSize - 1
		}
		log.Println("===start", start)
		log.Println("===end", end)
		log.Println("===diff", end-start)
		data, err := aliClient.downloadRange(url, start, end)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
				return
			}
		}
		iter += 1
		cnt += len(data)
		n, _ := buf.Write(data)

		log.Println("===current", iter)
		log.Println("===byte read", len(data))
		log.Println("===byte write", n)
		log.Println("===byte total", cnt)
	}

	log.Println("[success] result:")
	log.Printf("totalIter: %d, iter: %d\n", totalIter, iter)
	log.Printf("fileSize: %d, bytes send: %d\n", aliClient.fileSize, cnt)
	log.Printf("===sum %x\n", md5.Sum(buf.Bytes()))
	log.Println("===expect", checksum)
}

func TestAliClient_UploadChunk(t *testing.T) {
	reader := bytes.NewReader([]byte("测试内容2333测试内容2333测试内容2333"))
	siz := reader.Size()
	res := make([]byte, 0, siz)
	buf := bytes.NewBuffer(res)
	for {
		n, err := io.CopyN(buf, reader, siz)
		log.Println("===n", n)
		log.Println("===err", err)
		log.Println("===buf", buf.Bytes())
		if err != nil {
			break
		}
	}
	//for {
	//	n, err := reader.Read(buf)
	//	log.Println("===n", n)
	//	log.Println("===err", err)
	//	log.Println("===res", buf)
	//	if err != nil {
	//		if err == io.EOF {
	//			res = append(res, buf[:n]...)
	//			break
	//		}
	//		log.Println(err)
	//		return
	//	}
	//	res = append(res, buf[:n]...)
	//}
	log.Println("===finish", string(buf.Bytes()))
}

func TestDownloadFile(t *testing.T) {
	dir, _ := os.Getwd()
	filePath := filepath.Join(dir, "./3DMM_V12.safetensors")
	log.Println("===filePath", filePath)
	aliClient := New()
	body, _ := aliClient.bucket.GetObject("lora/COOLKIDS_MERGE_V2.5.safetensors")
	defer body.Close()
	//file, _ := os.Create(filePath)
	//io.Copy(file, body)
	bytes, _ := io.ReadAll(body)
	log.Println("===data", bytes)
}

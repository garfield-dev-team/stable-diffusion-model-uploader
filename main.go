package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"stable-diffusion-model-uploader/pkg/api"
	"stable-diffusion-model-uploader/pkg/client"
	"strings"
	"sync/atomic"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/sourcegraph/conc/pool"
)

var (
	totalCount  int
	uploadCount int64
	ignoreCount int64
	failedCount int64
	start       time.Time
)

func main() {
	log.Println("OSS Go SDK Version: ", oss.Version)

	log.Println("[info] getting model list...")
	list := api.FetchModelList()

	totalCount = len(list)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动 executor，用来监听并执行 worker goroutine 发送的任务
	executor := client.NewExecutor(ctx)
	executor.Run()

	p := pool.New().
		WithMaxGoroutines(runtime.NumCPU() * 1)

	start = time.Now()
	log.Println("[info] upload model to aliyun oss...")
	for _, item := range list {
		item := item
		p.Go(func() {
			aliClient := client.New()
			aliClient.UploadRange(item)
			err := aliClient.Error()
			if err != nil {
				if errors.Is(err, client.ErrObjectExist) {
					// 文件已存在
					atomic.AddInt64(&ignoreCount, 1)
					log.Println("[info] error object exist", item.Id)
				} else {
					// 文件上传失败
					atomic.AddInt64(&failedCount, 1)
					log.Println("[warn] upload failed", item.Id)
					log.Printf("%+v\n", err)
				}
				return
			}
			// 上传成功
			atomic.AddInt64(&uploadCount, 1)
			log.Println("[info] upload success", item.Id)
		})
	}

	p.Wait()

	// 上传任务完成
	// 主 goroutine 发送信号并等待 executor 退出
	cancel()
	<-executor.C

	var sb strings.Builder
	fmt.Fprintf(&sb, "[success] upload finished\n")
	fmt.Fprintf(&sb, "Total: %d\n", totalCount)
	fmt.Fprintf(&sb, "Upload: %d\n", uploadCount)
	fmt.Fprintf(&sb, "Failed: %d\n", failedCount)
	fmt.Fprintf(&sb, "Ignored: %d\n", ignoreCount)
	fmt.Fprintf(&sb, "Time: %.2fs", time.Since(start).Seconds())

	log.Println(sb.String())
}

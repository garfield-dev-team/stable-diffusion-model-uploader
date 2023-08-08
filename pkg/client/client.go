package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"

	"stable-diffusion-model-uploader/pkg/config"
	"stable-diffusion-model-uploader/pkg/model"
	"stable-diffusion-model-uploader/pkg/utils"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	client *oss.Client
	bucket *oss.Bucket
	once   sync.Once
)

var ErrObjectExist = fmt.Errorf("error object exists")

func InitOSS() {
	var err error
	client, err = oss.New(config.Endpoint, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		panic(fmt.Errorf("fail to init oss client: %w", err))
	}
	bucket, err = client.Bucket(config.BucketName)
	if err != nil {
		panic(fmt.Errorf("fail to get bucket: %w", err))
	}
}

func GetOSS() (*oss.Client, *oss.Bucket) {
	// 借助 sync.Once 实现懒汉式单例，在需要的时候初始化 SDK
	// 也可以在 init 函数初始化，缺点是会干扰单元测试
	once.Do(func() {
		InitOSS()
	})
	return client, bucket
}

type AliClient struct {
	client  *oss.Client
	bucket  *oss.Bucket
	nextPos int64
	buffer  []byte
	err     error
}

func New() *AliClient {
	client, bucket := GetOSS()
	// 缓冲区大小为 100MB
	buffer := make([]byte, 100*1024*1024)
	return &AliClient{
		client: client,
		bucket: bucket,
		buffer: buffer,
	}
}

func (c *AliClient) Error() error {
	return c.err
}

func (c *AliClient) UploadChunk(model *model.IModelDetailDTO) {
	resp, err := http.Get(model.DownloadUrl)
	if err != nil {
		c.err = fmt.Errorf("failed to download model: %w", err)
		return
	}
	defer resp.Body.Close()

	contentDisposition := resp.Header.Get("Content-Disposition")
	filename, err := utils.GetDownloadFileName(contentDisposition)
	if err != nil {
		c.err = err
		return
	}
	objectName := utils.GetObjectName(model.Type, filename)

	// 上传前判断文件是否存在
	exist, err := c.bucket.IsObjectExist(objectName)
	if err != nil {
		c.err = fmt.Errorf("failed to call IsObjectExist: %w", err)
		return
	}
	if exist {
		c.err = fmt.Errorf("%w", ErrObjectExist)
		return
	}
	option := []oss.Option{
		// 指定该Object被下载时的网页缓存行为。
		oss.CacheControl("no-cache"),
		// 指定该Object被下载时的名称。
		oss.ContentDisposition(fmt.Sprintf("attachment;filename=%s", objectName)),
		// 指定该Object的内容编码格式。
		oss.ContentEncoding("gzip"),
		// 指定Object的存储类型。
		oss.ObjectStorageClass(oss.StorageStandard),
		// 指定Object的访问权限。
		//oss.ObjectACL(oss.ACLPrivate),
		// 指定服务器端加密方式。
		//oss.ServerSideEncryption("AES256"),
		// 创建AppendObject时可以添加x-oss-meta-*，继续追加时不可以携带此参数。如果配置以x-oss-meta-*为前缀的参数，则该参数视为元数据。
		//oss.Meta("x-oss-meta-author", "Alice"),
	}
	for {
		// 从 HTTP 响应体中读取缓冲区大小的数据
		n, err := resp.Body.Read(c.buffer)
		if err != nil {
			if err == io.EOF {
				break // 读取完毕，退出循环
			}
			c.err = fmt.Errorf(
				"failed to read resp, objectName: %s, detail: %w",
				objectName, err)
			return
		}

		// 将缓冲区中的数据流上传到 OSS 上
		c.nextPos, err = c.bucket.AppendObject(objectName, bytes.NewReader(c.buffer[:n]), c.nextPos, option...)
		if err != nil {
			c.err = fmt.Errorf("failed to upload model to OSS, objectName: %s, detail: %w",
				objectName, err)
			return
		}
	}
}

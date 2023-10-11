package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"stable-diffusion-model-uploader/pkg/model"
)

func FetchModelList() []*model.IModelDetailDTO {
	url := "https://ke-api.study.163.com/ap/model"
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(fmt.Errorf("fail to create req: %w", err))
	}
	req.Header.Set("Cookie", "xxx")
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Errorf("fail to send req: %w", err))
	}
	defer resp.Body.Close()
	jsonData, err := io.ReadAll(resp.Body)
	var modelList *model.IModelListResp
	err = json.Unmarshal(jsonData, &modelList)
	if err != nil {
		panic(fmt.Errorf("fail to resolve json: %w", err))
	}

	log.Printf("%+v\n", modelList)
	if modelList.Code != 0 {
		panic(fmt.Errorf("fail to request: %s", modelList.Message))
	}

	l := len(modelList.Result.List)
	m := make(map[int]struct{}, l)
	// 注意切片预分配内存写法 make([]T, 0, len)
	res := make([]*model.IModelDetailDTO, 0, l)
	// 过滤 revAnimated_v122 模型
	m[10055] = struct{}{}

	for _, dto := range modelList.Result.List {
		if dto.Type != 0 {
			continue
		}
		if _, ok := m[dto.Id]; !ok {
			m[dto.Id] = struct{}{}
			res = append(res, dto)
		}
	}

	return res
}

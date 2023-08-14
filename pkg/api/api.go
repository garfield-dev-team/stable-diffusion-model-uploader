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
	req.Header.Set("Cookie", "HAS_VISIT_SD_WEBUI=true; EDUWEBDEVICE=61001fcbcb9f41e6a0317483f66dae9b; EDU-YKT-MODULE_GLOBAL_PRIVACY_DIALOG=true; OUTFOX_SEARCH_USER_ID_NCOO=307420402.20727146; hb_MA-BFF5-63705950A31C_source=sd.study.163.com; utm=\"eyJjIjoiIiwiY3QiOiIiLCJpIjoiIiwibSI6IiIsInMiOiIiLCJ0IjoiIn0=|aHR0cHM6Ly9zZC5zdHVkeS4xNjMuY29tLw==\"; STUDY_SESS=\"pzx2LuQgFzdPhSyxfC22K/ZfSVwfSQoaWPxZH3AwpPEYTBRB0PaKsPVbd/mH2Sh9lTS2z38pY/CwHu3jaiHvZwdKxsd+A4sQ8z6Zn74+ld3RmDNpMv76/n/T0nEjkAeXBjuH//koZcgBpm99y5p6REWrdTxX1zv/jpZ14OalGUILhur2Nm2wEb9HcEikV+3FTI8+lZKyHhiycNQo+g+/oA==\"; STUDY_INFO=\"oP4xHuH_Yhivb2un05qwMVbar_Qk|6|1031410010|1691992845104\"")
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

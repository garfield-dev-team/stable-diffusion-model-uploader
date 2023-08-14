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
	req.Header.Set("Cookie", "HAS_VISIT_SD_WEBUI=true; EDUWEBDEVICE=61001fcbcb9f41e6a0317483f66dae9b; EDU-YKT-MODULE_GLOBAL_PRIVACY_DIALOG=true; OUTFOX_SEARCH_USER_ID_NCOO=307420402.20727146; hb_MA-BFF5-63705950A31C_source=sd.study.163.com; NTESSTUDYSI=3cc26eec984a4f9bb214496034bcb7ea; utm=eyJjIjoiIiwiY3QiOiIiLCJpIjoiIiwibSI6IiIsInMiOiIiLCJ0IjoiIn0=|aHR0cHM6Ly9zdHVkeS4xNjMuY29tLw==; STUDY_SESS=\"pzx2LuQgFzdPhSyxfC22K/ZfSVwfSQoaWPxZH3AwpPEYTBRB0PaKsPVbd/mH2Sh9lTS2z38pY/CwHu3jaiHvZ3Ns7IVgVhVb3kWw90+U7+g2GezebVhZJ1iKgKbmBbur8eidZi8tyxWGICpjqp2QMUDTl4E8BVVA8n6AW7Yszk8Lhur2Nm2wEb9HcEikV+3FTI8+lZKyHhiycNQo+g+/oA==\"; STUDY_INFO=\"oP4xHuH_Yhivb2un05qwMVbar_Qk|6|1031410010|1691976511976\"")
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

	for _, dto := range modelList.Result.List {
		if _, ok := m[dto.Id]; !ok {
			m[dto.Id] = struct{}{}
			res = append(res, dto)
		}
	}

	return res
}

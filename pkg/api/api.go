package api

import (
	"encoding/json"
	"fmt"
	"io"
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
	req.Header.Set("Cookie", "HAS_VISIT_SD_WEBUI=true; EDUWEBDEVICE=61001fcbcb9f41e6a0317483f66dae9b; EDU-YKT-MODULE_GLOBAL_PRIVACY_DIALOG=true; OUTFOX_SEARCH_USER_ID_NCOO=307420402.20727146; NTESSTUDYSI=438aefd17ab54e46bf70740ea902ff7d; STUDY_SESS=\"pzx2LuQgFzdPhSyxfC22K/ZfSVwfSQoaWPxZH3AwpPEYTBRB0PaKsPVbd/mH2Sh9lTS2z38pY/CwHu3jaiHvZw79dNLUgKnb3xJu6RTESpq85qPJ/pA7pwIRhgxaYRcuBbG0VTifd7yjh7Jofer/abufEfg/rcCH9AKgD9S2xkYLhur2Nm2wEb9HcEikV+3FTI8+lZKyHhiycNQo+g+/oA==\"; STUDY_INFO=\"oP4xHuH_Yhivb2un05qwMVbar_Qk|6|1031410010|1691373584824\"; NETEASE_WDA_UID=1031410010#|#1639016325827; NTES_STUDY_YUNXIN_ACCID=s-1031410010; NTES_STUDY_YUNXIN_TOKEN=ab86cb4caa62f6162e07cde922b92ac8; hb_MA-BFF5-63705950A31C_source=sd.study.163.com")
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

	if modelList.Code != 0 {
		panic(fmt.Errorf("fail to request: %s", modelList.Message))
	}

	l := len(modelList.Result.List)
	m := make(map[int]struct{}, l)
	res := make([]*model.IModelDetailDTO, l)

	for _, dto := range modelList.Result.List {
		if _, ok := m[dto.Id]; !ok {
			m[dto.Id] = struct{}{}
			res = append(res, dto)
		}
	}

	return res
}

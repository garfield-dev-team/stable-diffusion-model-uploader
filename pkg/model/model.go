package model

type IModelDetailDTO struct {
	Id          int    `json:"id"`
	Type        uint8  `json:"type"` //checkpoint = 0,lora = 1,textual = 2
	DownloadUrl string `json:"downloadUrl"`
}

type IModelList struct {
	List []*IModelDetailDTO `json:"list"`
}

type IModelListResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  *IModelList `json:"result"`
}

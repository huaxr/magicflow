package playbookservice

import "github.com/huaxr/magicflow/component/dao/orm/models"

type CreatePlayBookReq struct {
	AppId       int    `json:"app_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreatePlayBookResp struct {
	PlayBookId int `json:"playbook_id"`
}

type SwitchVersionReq struct {
	PlayBookId int `json:"playbook_id"`
	SnapShotId int `json:"snapshot_id"`
}

type SwitchVersionResp struct {
	PlayBookId int `json:"playbook_id"`
}

type GetPlayBookReq struct {
	PlayBookId string `json:"playbook_id"`
}

type GetPlayBookResp struct {
	models.Playbook `json:"playbook"`
}

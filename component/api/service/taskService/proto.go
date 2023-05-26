// Author: huaxr
// Time:   2021/12/16 下午4:25
// Git:    huaxr

package taskService

type CreateTaskReq struct {
	Name    string `json:"name"`
	Service string `json:"service"` // proxy
	Region  string `json:"region"`  // cn
	// map[string]interface{} for config but receive from front-end it was string type
	Configuration string `json:"configuration"`
	AppId         int    `json:"app_id"`
	Description   string `json:"description"`
	InputExample  string `json:"input_example"`
	OutputExample string `json:"output_example"`
}

type CreateNormalTaskReq struct {
	CreateTaskReq
	TaskName string `json:"task_name"`
}

type CreateLocalTaskReq struct {
	CreateTaskReq
	PlaybookId int `json:"playbook_id"`
}

type CreateRemotePlaybookTaskReq struct {
	CreateTaskReq
	PlaybookId int    `json:"playbook_id"`
	RemoteApp  string `json:"remote_app"`
}

type CreateTaskRes struct {
	Id int `json:"id"`
}

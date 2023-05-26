// Author: huaxr
// Time:   2022/1/14 上午10:39
// Git:    huaxr

package request

type TicketResult struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Ticket  string `json:"ticket"`
}

type UserInfoResult struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Data    struct {
		Account  string `json:"account"`
		Name     string `json:"name"`
		Workcode string `json:"workcode"`
	} `json:"data"`
}

type UserMoreInfoResult struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Data    struct {
		TypeCode string `json:"type_code"`
		DeptInfo []struct {
			EhrDeptId    string `json:"ehr_dept_id"`
			DeptName     string `json:"dept_name"`
			DeptFullName string `json:"dept_full_name"`
		} `json:"dept_info"`
		Email  string `json:"email"`
		Avatar string `json:"avatar"`
	} `json:"data"`
}

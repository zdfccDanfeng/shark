package model

// Go语言中的 map 在并发情况下，只读是线程安全的，同时读写是线程不安全的。
var Task2Idp = map[string]string{}

type CustomUserProfileTag struct {
	Id                       int64  `json:"id" form:"id"`                                             // 任务id
	Table_Name               string `json:"table_name" form:"table_name"`                             // hive表名称
	Tag_Name                 string `json:"tag_name" form:"tag_name"`                                 // 标签名称
	Last_Udpate_Success_Time string `json:"last_udpate_success_time" form:"last_udpate_success_time"` // 0 正常状态， 1删除
	Update_status            int    `json:"update_status" form:"update_status"`                       // 更新状态
	Update_level             int    `json:"update_level" from "update_level"`                         // 优先级
	Data                     string `json:"data" from "data"`
	IdpTaskUrl               string `json:"idp_task_url"`
}

type Response struct {
	CODE int                    `json:"code"`
	MSG  string                 `json:"msg"`
	DATA []CustomUserProfileTag `json:"data"`
}

type GlobalResponse struct {
	CODE int         `json:"code"`
	MSG  string      `json:"msg"`
	DATA interface{} `json:"data"`
}

func QueryTaskList() []CustomUserProfileTag {
	return nil
}

type Message struct {
	Name  string `json:"name"`
	City  string `json:"city"`
	Other string `json:"other"`
}

type CustomResponse struct {
	Code int    `json:"code"` // 状态码
	Dsl  string `json:"dsl"`  // dsl语句
	Msg  string `json:"msg"`  // 返回信息
}

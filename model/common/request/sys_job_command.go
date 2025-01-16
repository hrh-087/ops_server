package request

type SysJobCommand struct {
	BatchType  int8  `json:"batchType" form:"batchType"`
	CommandId  uint  `json:"commandId" form:"commandId"`
	ServerList []int `json:"serverList" form:"serverList"`
}

package response

type PageResult struct {
	List     interface{} `json:"rows"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

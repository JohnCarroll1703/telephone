package schema

type Paginate struct {
	Records      int64 `json:"records"`
	Limit        int   `json:"limit"`
	Page         int   `json:"page"`
	TotalRecords int64 `json:"total_records"`
	TotalPage    int   `json:"total_page"`
}

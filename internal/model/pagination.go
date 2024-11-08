package model

type Paginate struct {
	Records      int64 `json:"records"`
	Limit        int   `json:"limit"`
	Page         int   `json:"page"`
	TotalRecords int64 `json:"total_records"`
	TotalPage    int   `json:"total_page"`
}

func (p *Paginate) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Paginate) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Paginate) GetLimit() int {
	if p.Limit <= 0 {
		p.Limit = 10
	}
	return p.Limit
}

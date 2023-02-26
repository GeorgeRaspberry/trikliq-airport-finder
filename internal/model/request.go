package model

type Request struct {
	Entity   string         `json:"entity"`
	Data     map[string]any `json:"data"`
	Metadata Metadata       `json:"metadata"`
}

type Metadata struct {
	Filter  Filter  `json:"filter"`
	OrderBy OrderBy `json:"orderBy"`
	//GroupBy string    `json:"groupBy"`
	//SumBy  SumBy    `json:"sumBy"`
	Fields []string `json:"fields"`
	Limit  int      `json:"limit"`
	Page   int      `json:"page"`
}

type Response struct {
	Status bool     `json:"status"`
	Errors []string `json:"errors"`
	Data   any      `json:"data,omitempty"`
	Limit  int      `json:"limit,omitempty"`
	Page   int      `json:"page,omitempty"`
}

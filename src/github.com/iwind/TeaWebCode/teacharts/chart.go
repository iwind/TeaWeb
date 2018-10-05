package teacharts

type ChartInterface interface {
	UniqueId() string
	SetUniqueId(id string)
}

type Chart struct {
	Id     string `json:"id"` // 用来标记图标的唯一性，可以不填，系统会自动生成
	Type   string `json:"type"`
	Name   string `json:"name"`
	Detail string `json:"detail"`
}

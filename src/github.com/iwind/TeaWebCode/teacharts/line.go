package teacharts

type Line struct {
	Name   string        `json:"name"`
	Values []interface{} `json:"values"`
}

type LineChart struct {
	Chart

	Lines  []*Line  `json:"lines"`
	Labels []string `json:"labels"`
}

func NewLineChart() *LineChart {
	p := &LineChart{}
	p.Type = "line"
	p.Lines = []*Line{}
	return p
}

func (this *LineChart) UniqueId() string {
	return this.Id
}

func (this *LineChart) SetUniqueId(id string) {
	this.Id = id
}

func (this *LineChart) AddLine(line *Line) {
	this.Lines = append(this.Lines, line)
}

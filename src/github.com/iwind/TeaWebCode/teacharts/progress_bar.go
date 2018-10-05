package teacharts

type ProgressBar struct {
	Chart

	Value float64 `json:"value"`
}

func NewProgressBar() *ProgressBar {
	p := &ProgressBar{}
	p.Type = "progressBar"
	return p
}

func (this *ProgressBar) UniqueId() string {
	return this.Id
}

func (this *ProgressBar) SetUniqueId(id string) {
	this.Id = id
}

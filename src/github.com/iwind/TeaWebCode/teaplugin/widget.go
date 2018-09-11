package teaplugin

type Widget struct {
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	MoreURL   string `json:"moreUrl"`
	TopBar    bool   `json:"topBar"`
	MenuBar   bool   `json:"menuBar"`
	HelperBar bool   `json:"helperBar"`
	Dashboard bool   `json:"dashboard"`
}

func NewWidget() *Widget {
	return &Widget{}
}

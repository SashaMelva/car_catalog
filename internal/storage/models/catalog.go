package model

type CarCatalog struct {
	Cars []*Car `json:"cars"`
}

type RegNumsCatalog struct {
	RegNums []string `json:"regNums"`
}

type RequestBody struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Content     *Car   `json:"content"`
}

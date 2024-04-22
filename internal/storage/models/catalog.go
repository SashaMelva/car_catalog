package model

type CarCatalog struct {
	Cars []*Car `json:"cars"`
}

type RegNumsCatalog struct {
	RegNums []string `json:"regNums"`
}

package model

type Car struct {
	RegNum string  `json:"regNum"  db:"reg_num"`
	Mark   string  `json:"mark"    db:"mark"`
	Model  string  `json:"model"   db:"model"`
	Year   int     `json:"year"    db:"year"`
	Owner  *People `json:"owner"   db:"owner"`
}

package models

type BaseSite struct {
	Title   string   `json:"title"`
	Content any      `json:"body"`
	Footer  []string `json:"footer"`
}

type ArticleList struct {
	Article Article `json:"articles"`
	List    any     `json:"list"`
}

type ErrorSite struct {
	Message string
	Status  int
}

package ingestor

type Document struct {
	ID      int64  `json:"id"`
	Text    string `json:"text"`
	Source  string `json:"source"`
	Content string `json:"content"`
}

package models

type Video struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
	URL   string   `json:"url"`
	Src   string   `json:"src"`
	// Avatar   string   `json:"avatar"`
	Username string `json:"username"`
}

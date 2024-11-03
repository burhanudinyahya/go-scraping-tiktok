package models

type Video struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Tags     []string `json:"tags"`
	URL      string   `json:"url"`
	Src      string   `json:"src"`
	Image    string   `json:"image"`
	Username string   `json:"username"`
	Cookie   string   `json:"cookie"`
	Likes    string   `json:"likes"`
	Comments string   `json:"comments"`
}

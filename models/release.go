package models

// Release represents release information from Github
type Release struct {
	Name        string `json:"name"`
	TagName     string `json:"tag_name"`
	URL         string `json:"html_url"`
	Body        string `json:"body"`
	CreatedAt   string `json:"created_at"`
	PublishedAt string `json:"published_at"`
}

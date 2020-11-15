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

// ReleaseWorker starts worker to get project latest release.
func ReleaseWorker(jobs <-chan Project, results chan<- Release) {
	for project := range jobs {
		release, err := project.GetInformation()
		if err == nil {
			results <- release
		}
	}
}

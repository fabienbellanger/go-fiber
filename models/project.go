package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

// Project represents a Github project.
type Project struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// GetInformation calls Github API to access last release information.
func (p *Project) GetInformation() (Release, error) {
	url := fmt.Sprintf("%s/repos/%s/releases/latest", viper.GetString("github.apiBaseURL"), p.URL)

	var release Release

	// Requête vers Github
	// -------------------
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return release, err
	}
	req.SetBasicAuth(viper.GetString("github.apiUsername"), viper.GetString("github.apiToken"))
	resp, err := client.Do(req)
	if err != nil {
		return release, err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return release, err
	}

	// Récupération de la dernière release
	// -----------------------------------
	err = json.Unmarshal(bodyText, &release)
	if err != nil {
		return release, err
	}
	return release, nil
}

// LoadProjectsFromFile loads projects from JSON file.
func LoadProjectsFromFile(file string) ([]Project, error) {
	projects := make([]Project, 0)

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return projects, err
	}

	err = json.Unmarshal(content, &projects)
	if err != nil {
		return projects, err
	}
	return projects, nil
}

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

// New initializes a new project.
func (p *Project) New(name, url string) {
	p.Name = name
	p.URL = url
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

	// Recupération de la dernière release
	// -----------------------------------
	err = json.Unmarshal(bodyText, &release)
	if err != nil {
		return release, err
	}
	return release, nil
}

package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// Project represents a Github project.
type Project struct {
	Name string `json:"name"`
	Repo string `json:"repo"`
}

// Release represents release information from Github
type Release struct {
	Project     Project `json:"project"`
	Name        string  `json:"name"`
	TagName     string  `json:"tag_name"`
	URL         string  `json:"html_url"`
	Body        string  `json:"body"`
	CreatedAt   string  `json:"created_at"`
	PublishedAt string  `json:"published_at"`
}

type releasesCache struct {
	releases []Release
	expireAt time.Time
	mux      sync.RWMutex
}

var (
	// Used to cached Github response
	cachedReleases releasesCache = releasesCache{}
)

// GetInformation calls Github API to access last release information.
func (p *Project) GetInformation() (Release, error) {
	url := fmt.Sprintf("%s/repos/%s/releases/latest", viper.GetString("github.apiBaseURL"), p.Repo)

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

	// Ajout information sur le projet
	// -------------------------------
	release.Project = *p

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

// releaseWorker starts worker to get project latest release.
func releaseWorker(jobs <-chan Project, results chan<- Release) {
	for project := range jobs {
		release, err := project.GetInformation()
		if err == nil {
			results <- release
		}
	}
}

// getLatestReleases find all projects latest release.
func getLatestReleases(projects []Project) ([]Release, error) {
	numProjects := len(projects)
	jobs := make(chan Project, numProjects)
	results := make(chan Release, numProjects)

	// Nombre de workers
	// -----------------
	numWorkers := runtime.NumCPU()
	if numProjects < numWorkers {
		numWorkers = numProjects
	}

	// Lancement des workers
	// ---------------------
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			releaseWorker(jobs, results)
		}()
	}

	// Fermeture du channel results quand tous les workers ont terminés
	// ----------------------------------------------------------------
	go func() {
		defer close(results)
		wg.Wait()
	}()

	// Envoi des jobs
	// --------------
	go func() {
		defer close(jobs)
		for _, project := range projects {
			jobs <- project
		}
	}()

	// Traitement des résultats
	// ------------------------
	releases := make([]Release, 0)
	for r := range results {
		releases = append(releases, r)
	}

	return releases, nil
}

// ReleasesProcess returns latest releases.
func ReleasesProcess(projects []Project) ([]Release, error) {
	cachedReleases.mux.Lock()
	defer cachedReleases.mux.Unlock()
	now := time.Now()
	if len(cachedReleases.releases) == 0 || cachedReleases.expireAt.Before(now) {
		r, err := getLatestReleases(projects)
		if err != nil {
			return r, err
		}

		cachedReleases.releases = r
		cachedReleases.expireAt = now.Local().Add(time.Hour)
	}
	return cachedReleases.releases, nil
}

package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// Project represents a Github project.
type Project struct {
	Name     string `json:"name"`
	Repo     string `json:"repo"`
	Language string `json:"language"`
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
	projects []Project
	releases []Release
	expireAt time.Time
	mux      sync.RWMutex
}

var (
	// Used to cached Github response
	cachedReleases releasesCache = releasesCache{}
)

func (p Project) String() string {
	return fmt.Sprintf("[%s] %s (%s)", p.Language, p.Name, p.Repo)
}

// info calls Github API to access last release information.
func (p *Project) info() (Release, error) {
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

	// Si status code != 200, on retourne une erreur
	// ---------------------------------------------
	if resp.StatusCode != http.StatusOK {
		return release, fmt.Errorf("error during retrieving Github %s project: status code=%d", *p, resp.StatusCode)
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

// loadProjectsFromFile loads projects from JSON file.
func loadProjectsFromFile(file string) ([]Project, error) {
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
		release, err := project.info()
		if err == nil {
			results <- release
		} else {
			log.Printf("%v\n", err)
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

// GetReleases returns latest release of Github projects.
func GetReleases() ([]Release, time.Time, error) {
	cachedReleases.mux.Lock()
	defer cachedReleases.mux.Unlock()

	now := time.Now()
	if len(cachedReleases.releases) == 0 || cachedReleases.expireAt.Before(now) {
		projects, err := loadProjectsFromFile("projects.json")
		log.Printf("Projects: %v\n", projects)
		if err != nil {
			return cachedReleases.releases, time.Now(), err
		}

		r, err := getLatestReleases(projects)
		if err != nil {
			return r, time.Now(), err
		}

		cachedReleases.projects = projects
		cachedReleases.releases = r
		cachedReleases.expireAt = now.Local().Add(time.Hour)
	}
	return cachedReleases.releases, cachedReleases.expireAt, nil
}

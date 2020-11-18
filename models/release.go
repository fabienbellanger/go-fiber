package models

import (
	"runtime"
	"sync"
	"time"
)

// Release represents release information from Github
type Release struct {
	Name        string `json:"name"`
	TagName     string `json:"tag_name"`
	URL         string `json:"html_url"`
	Body        string `json:"body"`
	CreatedAt   string `json:"created_at"`
	PublishedAt string `json:"published_at"`
}

type ReleasesCache struct {
	Releases []Release
	ExpireAt time.Time
	Mux      sync.RWMutex
}

var CachedReleases ReleasesCache = ReleasesCache{}

// releaseWorker starts worker to get project latest release.
func releaseWorker(jobs <-chan Project, results chan<- Release) {
	for project := range jobs {
		release, err := project.GetInformation()
		if err == nil {
			results <- release
		}
	}
}

// ReleasesProcess find all projects latest release.
func ReleasesProcess(projects []Project) ([]Release, error) {
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

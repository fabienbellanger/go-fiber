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

type releasesCache struct {
	releases []Release
	expireAt time.Time
	mux      sync.RWMutex
}

var (
	// Used to cached Github response
	cachedReleases releasesCache = releasesCache{}
)

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

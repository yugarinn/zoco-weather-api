package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const CACHE_LIFETIME_IN_MINUTES = 5

func CacheCragsResponse(crags CragsList) {
	cleanCache()

	datetime := time.Now().UTC().Format(time.RFC3339)
	cacheFilename := fmt.Sprintf("/cache/%s-cached-crags.json", datetime)
	cacheFilepath := filepath.Join(os.Getenv("DATA_PATH"), cacheFilename)
    f, err := os.Create(cacheFilepath)

    if err != nil {
        panic(err)
    }

    defer f.Close()

    encoder := json.NewEncoder(f)
    err = encoder.Encode(&crags)

    if err != nil {
        panic(err)
    }
}

func GetCachedCrags() (list CragsList, success bool) {
	cacheFoldername := fmt.Sprintf("%s/cache", os.Getenv("DATA_PATH"))
	caches, readDirErr := ioutil.ReadDir(cacheFoldername)

    if readDirErr != nil || len(caches) == 0 {
        return CragsList{}, false
    }

	cachedFilename := caches[len(caches) - 1].Name()

	if fileHasExpired(cachedFilename) {
        return CragsList{}, false
	}

	cacheFilename := fmt.Sprintf("%s/cache/%s", os.Getenv("DATA_PATH"), cachedFilename)
	cachedFile, _ := ioutil.ReadFile(cacheFilename)
	crags := CragsList{}

	unmarshalErr := json.Unmarshal([]byte(cachedFile), &crags)

    if unmarshalErr != nil {
        return CragsList{}, false
    }

	return crags, true
}

func fileHasExpired(cachedFilename string) bool {
    parts := strings.Split(cachedFilename, "Z")

    if len(parts) < 2 {
        return true
    }

    cachedAt, err := time.Parse("2006-01-02T15:04:05", parts[0])

    if err != nil {
		fmt.Println(err)
        return true
    }

    now := time.Now().UTC()
    diff := now.Sub(cachedAt)

    if diff.Minutes() > CACHE_LIFETIME_IN_MINUTES {
        return true
    } else {
        return false
    }
}

func cleanCache() {
	cacheFoldername := fmt.Sprintf("%s/cache", os.Getenv("DATA_PATH"))
	caches, readDirErr := ioutil.ReadDir(cacheFoldername)

    if readDirErr != nil {
        return
    }

    for _, file := range caches {
        err := os.Remove(cacheFoldername + "/" + file.Name())

        if err != nil {
            return
        }
    }
}

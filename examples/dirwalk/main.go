package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/dbaggerman/cuba"
)

type DirectoryJob struct {
	path string
	info os.FileInfo
}

func worker(item interface{}) []interface{} {
	job := item.(*DirectoryJob)

	file, err := os.Open(job.path)
	if err != nil {
		log.Printf("[ERR] Failed to open %s: %v", job.path, err)
		return nil
	}
	defer file.Close()

	var newJobs []interface{}

	dirents, err := file.Readdir(-1)
	if err != nil {
		log.Printf("[ERR] Failed to read %s dirnames: %v", job.path, err)
		return nil
	}

	for _, dirent := range dirents {
		direntPath := filepath.Join(job.path, dirent.Name())

		if !dirent.IsDir() {
			log.Printf("[FILE] %s", direntPath)
			continue
		} else {
			log.Printf("[DIR] %s", direntPath)
			direntJob := &DirectoryJob{
				path: direntPath,
				info: dirent,
			}
			newJobs = append(newJobs, direntJob)
		}
	}

	return newJobs
}

func main() {
	ws := cuba.NewStack(worker)

	info, err := os.Stat(".")
	if err != nil {
		panic(err)
	}

	root := &DirectoryJob{
		path: ".",
		info: info,
	}

	ws.Push([]interface{}{root})

	ws.Run()
}

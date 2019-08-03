package main

import (
	"log"
	"os"
	"path/filepath"
	// "strings"

	"github.com/dbaggerman/cuba"
	// "github.com/pkg/profile"
)

type DirectoryJob struct {
	path string
	info os.FileInfo
}

func worker(item interface{}) []interface{} {
	job := item.(*DirectoryJob)
	// log.Printf("[JOB] %s", job.path)


	file, err := os.Open(job.path)
	if err != nil {
		log.Printf("[ERR] Failed to open %s: %v", job.path, err)
		return nil
	}
	defer file.Close()

	// fileInfo, err := file.Stat()
	// if err != nil {
	// 	log.Printf("[ERR] Failed to stat %s: %v", job.path, err)
	// 	return nil
	// }

	// if !fileInfo.IsDir() {
	// 	log.Printf("[FILE] %s", job.path)
	// 	return nil
	// }

	// log.Printf("[DIR] %s", job.path)

	var newJobs []interface{}

	dirents, err := file.Readdir(-1)
	if err != nil {
		log.Printf("[ERR] Failed to read %s dirnames: %v", job.path, err)
		return nil
	}

	// haveIgnore := false

	for _, dirent := range dirents {
		// if name == ".ignore" {
		// 	haveIgnore = true
		// }

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

	// if haveIgnore {
	// 	log.Printf("Found ignore file: %s", filepath.Join(path, ".ignore"))
	// }

	return newJobs
}

func main() {
	// defer profile.Start().Stop()

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

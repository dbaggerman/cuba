package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/dbaggerman/cuba"
)

type Directory struct {
	path string
	info os.FileInfo
}

func worker(handle *cuba.Handle) {
	job := handle.Item().(*Directory)

	file, err := os.Open(job.path)
	if err != nil {
		log.Printf("[ERR] Failed to open %s: %v", job.path, err)
		return
	}
	defer file.Close()

	dirents, err := file.Readdir(-1)
	if err != nil {
		log.Printf("[ERR] Failed to read %s dirnames: %v", job.path, err)
		return
	}

	for _, dirent := range dirents {
		direntPath := filepath.Join(job.path, dirent.Name())

		if dirent.Name() == ".git" {
			continue
		}

		if !dirent.IsDir() {
			log.Printf("[FILE] %s", direntPath)
			continue
		} else {
			log.Printf("[DIR] %s", direntPath)
			handle.Push(
				&Directory{
					path: direntPath,
					info: dirent,
				},
			)
		}
	}
}

func main() {
	ws := cuba.New(worker, cuba.NewStack())

	info, err := os.Stat(".")
	if err != nil {
		panic(err)
	}

	root := &Directory{
		path: ".",
		info: info,
	}

	ws.Push(root)

	ws.Finish()
}

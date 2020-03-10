package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gosimple/slug"
)

func main() {
	extensions := map[string]bool{
		"":     false,
		".mkv": true,
		".avi": true,
		".m4v": true,
		".mp4": true,
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) == 2 {
		dir = os.Args[1]
	}

	dir, err = filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if extensions[filepath.Ext(path)] {
			re := regexp.MustCompile(`(?:^[A-Z])?(?:.*)(?:\/|\\\\?)([a-zA-Z ]?.*)(?:s|S|Season) ?([0-9]+) ?(?:e|E|Episode) ?([0-9]+)(?:.*)(?:\.)`)
			for _, loc := range re.FindAllStringSubmatch(path, -1) {
				slugifiedFileName := slug.Make(loc[1])

				if len(loc) == 4 && "" != slugifiedFileName {
					filename := fmt.Sprintf("%s_s%02se%02s", slugifiedFileName, loc[2], loc[3])
					newdir := fmt.Sprintf("%s/%s", dir, slugifiedFileName)

					if _, err := os.Stat(newdir); os.IsNotExist(err) {
						fmt.Printf("Directory \"%s\" doesn't exists. Creating it ... \n", newdir)
						err = os.Mkdir(newdir, os.ModePerm)
						if err != nil {
							log.Fatal(err)
						}
					}

					fmt.Printf("Moving file \"%s\" ... \n", path)
					err = os.Rename(path, fmt.Sprintf("%s/%s%s", newdir, filename, filepath.Ext(path)))
					if err != nil {
						log.Fatal(err)
					}
				} else {
					fmt.Println(loc)
				}
			}
		}

		return nil
	})
}

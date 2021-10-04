package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gosimple/slug"
)

var extensions = map[string]bool{
	"":     false,
	".mkv": true,
	".avi": true,
	".m4v": true,
	".mp4": true,
}

var regex = `(?:^[A-Z])?(?:.*)(?:\/|\\\\?)([a-zA-Z ]?.*)(?:s|S|Season) ?([0-9]+) ?(?:e|E|Episode) ?([0-9]+)(?:.*)(?:\.)`

func main() {
	var dir string
	flag.StringVar(&dir, "d", "", "The path in which the recursive search of media takes place")
	flag.Parse()

	if dir == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if extensions[filepath.Ext(path)] {
			source := filepath.Dir(path)
			base := filepath.Base(source)

			for _, loc := range regexp.MustCompile(regex).FindAllStringSubmatch(path, -1) {
				slugName := slug.Make(loc[1])
				if slugName == base {
					continue
				}

				if len(loc) == 4 && slugName != "" {
					filename := fmt.Sprintf("%s_s%02se%02s", slugName, loc[2], loc[3])
					newdir := filepath.Clean(fmt.Sprintf("%s/%s", dir, slugName))

					if _, err := os.Stat(newdir); os.IsNotExist(err) {
						err = os.Mkdir(newdir, os.ModePerm)
						if err != nil {
							log.Fatal(err)
						}
					}

					episode := fmt.Sprintf("%s%s", filename, filepath.Ext(path))
					fmt.Printf("[+] %s -> %s\n", episode, newdir)
					err = os.Rename(path, fmt.Sprintf("%s/%s", newdir, episode))
					if err != nil {
						log.Fatal(err)
					}

					// todo: remove directory from which the media was taken if its empty
				}
			}
		}

		return nil
	})
}

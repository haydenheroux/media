package main

import (
	"downloader/downloader"
	"downloader/track"
	"fmt"

	"flag"
	"log"
	"os"
)

const (
	APP_NAME                 = "downloader"
	DEFAULT_OUTPUT_DIRECTORY = ""
)

var (
	downloaderName  string
	outputFormat    string
	outputDirectory string
	printInfo       bool
)

func init() {
	flag.StringVar(&downloaderName, "d", "ytdl", "downloader name")
	flag.StringVar(&outputFormat, "f", "mp3", "output format")
	flag.StringVar(&outputDirectory, "o", DEFAULT_OUTPUT_DIRECTORY, "output directory")

	flag.BoolVar(&printInfo, "p", false, "print information as a track is downloading")
}

func main() {
	flag.Parse()

	logger := log.New(os.Stderr, APP_NAME+": ", 0)

	files := flag.Args()

	if len(files) == 0 {
		logger.Fatalf("no input files provided\n")
	}

	dl := downloader.CreateDownloader(downloaderName, outputFormat)

	if shouldMkdir() {
		if err := mkdir(outputDirectory); err != nil {
			logger.Fatalf("failed to make output directory (%v)\n", err)
		}
	}

	tracks, err := parseFiles(files)

	if err != nil {
		logger.Fatalf("failed to parse input file (%v)", err)
	}

	existing := scanExisting(tracks, dl, outputDirectory)

	if printInfo {
		for track := range existing {
			fmt.Printf("found: %s\n", track)
		}
	}

	tracks = removeExisting(tracks, existing)

	for _, track := range tracks {
		if printInfo {
			fmt.Printf("started: %s\n", track)
		}

		err := dl.Download(track, outputDirectory)
		if err != nil {
			logger.Fatalf("failed to download %s (%v)\n", track, err)
		}

		if printInfo {
			fmt.Printf("completed: %s\n", track)
		}
	}
}

func shouldMkdir() bool {
	return outputDirectory != DEFAULT_OUTPUT_DIRECTORY
}

func parseFiles(names []string) ([]track.Track, error) {
	result := make([]track.Track, 0)

	for _, name := range names {
		tracks, err := parseFile(name)

		if err != nil {
			return result, err
		}

		for _, track := range tracks {
			result = append(result, track)
		}
	}

	return result, nil
}

func parseFile(name string) ([]track.Track, error) {
	file, err := os.Open(name)
	defer file.Close()

	if err != nil {
		return []track.Track{}, err
	}

	tracks, err := track.ParseFile(file)
	if err != nil {
		return []track.Track{}, err
	}

	return tracks, nil
}

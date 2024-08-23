package main

import (
	"flag"
	"github.com/neveldo/komoot-gpx-exporter/internal/application/query"
	"github.com/neveldo/komoot-gpx-exporter/internal/infrastructure/file"
	"github.com/neveldo/komoot-gpx-exporter/internal/infrastructure/http"
	"log"
)

var (
	maxTours       = 10000
	maxParallelism = 10
)

func main() {
	var cookie = flag.String("cookie", "", "Komoot cookies header (you can retrieve it by browsing Komoot app in a web browser (mandatory)")
	var rootDir = flag.String("dir", ".", "Directory where to save GPX files (mandatory)")
	var userId = flag.String("user", "", "Komoot user ID. It appears in the URL (mandatory)")
	var sport = flag.String("sport", "", "Filter tours for a specific sport among hike, touringbicycle, mtb & jogging (optional)")

	flag.Parse()

	log.Printf("Running Komoot GPX exporter for user id %s into directory %s", *userId, *rootDir)

	komootHttpRepository := http.NewKomootHttpRepository(*cookie)
	gpxFileRepository := file.NewGpxFileRepository(*rootDir)

	handler := query.NewDownloadKomootToursHandler(komootHttpRepository, gpxFileRepository)

	handlerQuery := query.DownloadKomootTours{
		MaxTours:       maxTours,
		MaxParallelism: maxParallelism,
		UserId:         *userId,
		Sport:          *sport,
	}

	err := handler.Handle(handlerQuery)

	if err != nil {
		log.Fatal(err)
	}
}

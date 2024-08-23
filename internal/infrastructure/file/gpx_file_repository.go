package file

import (
	"fmt"
	"github.com/neveldo/komoot-gpx-exporter/internal/domain"
	"log"
	"os"
	"strings"
)

type GpxFileRepository struct {
	directory string
}

// SaveGPX saves the GPX data in a file
// The name of the file uses tour metadata date + title + ID
func (repository *GpxFileRepository) SaveGPX(tour domain.Tour, gpx []byte) error {
	filepath := fmt.Sprintf(
		"%s/%s - %s - %s.gpx",
		repository.directory,
		tour.Date[:10],
		strings.ReplaceAll(tour.Name, "/", "-"),
		tour.ID,
	)

	f, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("error %s", err)
	}
	defer f.Close()

	_, err = f.Write(gpx)
	if err != nil {
		log.Fatalf("error %s", err)
	}
	return nil
}

func NewGpxFileRepository(rootDir string) *GpxFileRepository {
	return &GpxFileRepository{directory: rootDir}
}

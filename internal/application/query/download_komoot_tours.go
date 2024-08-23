package query

import (
	"github.com/neveldo/komoot-gpx-exporter/internal/domain"
	"log"
	"sync"
)

type DownloadKomootToursHandler struct {
	komootRepository domain.KomootRepository
	gpxRepository    domain.GpxRepository
}

type DownloadKomootTours struct {
	MaxTours       int
	MaxParallelism int
	UserId         string
	Sport          string
}

// Handle handles DownloadKomootTours queries
// It downloads all GPX for the Komoot account specified by the UserId parameter.
// MaxTours parameter allows to specify the max mumber of tour to retrieve
// MaxParallelism parameter set the number of go routines to run
func (handler *DownloadKomootToursHandler) Handle(downloadKomootTours DownloadKomootTours) error {
	tours, err := handler.komootRepository.GetTours(downloadKomootTours.UserId, downloadKomootTours.Sport, downloadKomootTours.MaxTours)

	if err != nil {
		return err
	}

	toursChannel := make(chan domain.Tour)

	go feedToursChannel(tours, toursChannel)

	handler.downloadAllGpx(downloadKomootTours.MaxParallelism, toursChannel)

	return nil
}

// downloadAllGpx downloads all GPX by consuming values from the channel toursChannel
// maxParallelism parameter set the number of go routines to run
func (handler *DownloadKomootToursHandler) downloadAllGpx(maxParallelism int, toursChannel <-chan domain.Tour) {
	wg := sync.WaitGroup{}
	wg.Add(maxParallelism)
	for i := 0; i < maxParallelism; i++ {
		go func() {
			for tour := range toursChannel {
				log.Printf("Downloading tour \"%s\"", tour.Name)

				gpx, _ := handler.komootRepository.GetGPX(tour.ID)
				err := handler.gpxRepository.SaveGPX(tour, gpx)

				if err != nil {
					log.Print(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

// feedToursChannel sends each tour from the parameter tours into the channel toursChannel
func feedToursChannel(tours []domain.Tour, toursChannel chan<- domain.Tour) {
	for _, tour := range tours {
		toursChannel <- tour
	}
	close(toursChannel)
}

func NewDownloadKomootToursHandler(komootRepository domain.KomootRepository, gpxRepository domain.GpxRepository) *DownloadKomootToursHandler {
	return &DownloadKomootToursHandler{
		komootRepository: komootRepository,
		gpxRepository:    gpxRepository,
	}
}

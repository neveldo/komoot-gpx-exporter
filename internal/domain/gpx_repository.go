package domain

// GpxRepository defines an interface to save a GPX and tour metadata somewhere
type GpxRepository interface {
	SaveGPX(tour Tour, gpx []byte) error
}

package domain

// KomootRepository defines an interface to get all tours from a Komoot account
// And get a specific GPX track by its ID
type KomootRepository interface {
	GetTours(userId, sport string, limit int) ([]Tour, error)
	GetGPX(tourID string) ([]byte, error)
}

// Tour represents the structure for each tour with only the id and name fields.
type Tour struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Date string `json:"date"`
}

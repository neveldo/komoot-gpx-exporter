package http

import (
	"encoding/json"
	"errors"
	"github.com/neveldo/komoot-gpx-exporter/internal/domain"
	"io"
	"net/http"
	"strconv"
)

type KomootHttpRepository struct {
	cookie     string
	httpClient *http.Client
}

// embedded represents the _embedded field containing the tours array.
type embedded struct {
	Tours []domain.Tour `json:"tours"`
}

// apiResponse represents the top-level structure of the JSON response.
type apiResponse struct {
	Embedded embedded `json:"_embedded"`
}

// GetTours returns a list of Tour structs of all Komoot tours for the specified userId and
// up to the specified limit
func (komootClient *KomootHttpRepository) GetTours(userId, sport string, limit int) ([]domain.Tour, error) {
	url := "https://www.komoot.com/api/v007/users/" + userId + "/tours/?sport_types=" + sport + "&type=tour_recorded&sort_field=date&sort_direction=desc&name=&status=private&hl=fr&page=0&limit=" + strconv.Itoa(limit)

	// Create a new request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []domain.Tour{}, err
	}

	komootClient.setHeaders(request)

	// Send the request
	resp, err := komootClient.httpClient.Do(request)
	if err != nil {
		return []domain.Tour{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []domain.Tour{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []domain.Tour{}, errors.New("Error " + strconv.Itoa(resp.StatusCode) + " : " + resp.Status)
	}

	var response apiResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return []domain.Tour{}, err
	}
	return response.Embedded.Tours, nil
}

// GetGPX return the GPX track of a specific tour identified by tourID
func (komootClient *KomootHttpRepository) GetGPX(tourID string) ([]byte, error) {
	url := "https://www.komoot.com/tour/" + tourID + "/download"

	// Create a new request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}

	komootClient.setHeaders(request)

	// Send the request
	resp, err := komootClient.httpClient.Do(request)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func (komootClient *KomootHttpRepository) setHeaders(req *http.Request) {
	// Set the headers
	req.Header.Set("accept", "application/hal+json,application/json")
	req.Header.Set("accept-language", "fr")
	req.Header.Set("cookie", komootClient.cookie)
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
}

func NewKomootHttpRepository(cookie string) *KomootHttpRepository {
	return &KomootHttpRepository{
		cookie:     cookie,
		httpClient: &http.Client{},
	}
}

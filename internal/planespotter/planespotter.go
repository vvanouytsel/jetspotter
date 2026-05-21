package planespotter

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const userAgent = "jetspotter/1.0 (+https://github.com/vvanouytsel/jetspotter)"

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Thumbnail struct {
	Src  string `json:"src"`
	Size Size   `json:"size"`
}

type Image struct {
	ID             string    `json:"id"`
	Thumbnail      Thumbnail `json:"thumbnail"`
	ThumbnailLarge Thumbnail `json:"thumbnail_large"`
	Link           string    `json:"link"`
	Photographer   string    `json:"photographer"`
}

type ImagesData struct {
	Images []Image `json:"photos"`
}

// GetImageFromAPI uses the planespotters.net API to retrieve information about an image based on ICAO code.
func GetImageFromAPI(ICAO, registration string) (image *Image) {
	image = getImageByICAO(ICAO)
	if image == nil {
		image = getImageByRegistration(registration)
	}

	if image == nil {
		return &Image{}
	}

	return image
}

func getImageByICAO(ICAO string) (image *Image) {
	return fetchImage(fmt.Sprintf("https://api.planespotters.net/pub/photos/hex/%s", ICAO))
}

func getImageByRegistration(registration string) (image *Image) {
	return fetchImage(fmt.Sprintf("https://api.planespotters.net/pub/photos/reg/%s", registration))
}

func fetchImage(URL string) *Image {
	var images ImagesData
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Printf("planespotter: failed to create request for %s: %v", URL, err)
		return nil
	}
	req.Header.Set("User-Agent", userAgent)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("planespotter: request failed for %s: %v", URL, err)
		return nil
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		log.Printf("planespotter: unexpected status %d for %s: %s", res.StatusCode, URL, string(body))
		return nil
	}

	if err := json.Unmarshal(body, &images); err != nil {
		log.Printf("planespotter: failed to parse response from %s: %v", URL, err)
		return nil
	}

	if len(images.Images) == 0 {
		return nil
	}

	return &images.Images[0]
}

package planespotter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
	var images ImagesData
	URL := fmt.Sprintf("https://api.planespotters.net/pub/photos/hex/%s", ICAO)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &images)
	if err != nil {
		return nil
	}

	if len(images.Images) == 0 {
		return nil
	}

	return &images.Images[0]
}

func getImageByRegistration(registration string) (image *Image) {
	var images ImagesData
	URL := fmt.Sprintf("https://api.planespotters.net/pub/photos/reg/%s", registration)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &images)
	if err != nil {
		return nil
	}

	if len(images.Images) == 0 {
		return nil
	}

	return &images.Images[0]
}

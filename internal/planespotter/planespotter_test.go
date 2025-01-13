package planespotter

import "testing"

func TestGetImageFromAPI(t *testing.T) {
	// 44c1e5
	expectedThumbnailURL := "https://t.plnspttrs.net/21151/1449607_d9b2824f6f_280.jpg"
	expectedImageURL := "https://www.planespotters.net/photo/1449607/g-10-belgium-politie-police-mcdonnell-douglas-md-900-explorer?utm_source=api"
	image := GetImageFromAPI("44c1e5", "notimportant")

	if expectedThumbnailURL != image.ThumbnailLarge.Src {
		t.Fatalf("expected '%v' to be the same as '%v'", expectedThumbnailURL, image.ThumbnailLarge.Src)
	}

	if expectedImageURL != image.Link {
		t.Fatalf("expected '%v' to be the same as '%v'", expectedImageURL, image.Link)
	}
}

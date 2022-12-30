package spotify

import (
	"github.com/zmb3/spotify/v2"
	"strings"
)

func GetArtistNames(item *spotify.FullTrack) string {
	artistNames := make([]string, 0)
	for _, artist := range item.Artists {
		artistNames = append(artistNames, artist.Name)
	}
	return strings.Join(artistNames, " ")
}

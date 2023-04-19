package downloader

import "music_dl/track"

// Downloader is implemented by any type that is used to download tracks.
type Downloader interface {
	// Download performs the process of downloading the track.
	Download(track.Track) error
}
